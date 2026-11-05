---
name: scenes-and-nodes
description: >-
  Godot 4.7 scene/node architecture for GDScript and C#. Use when building the
  node tree, ordering _ready/_enter_tree logic, accessing child nodes, instancing
  scenes (PackedScene) at runtime, or freeing nodes — not for signals, physics or UI.
---

# Scenes & Nodes (Godot 4.7)

Nodes are the fundamental building blocks of a Godot game; organized into a tree with a single root, they
form a **scene**, and a game is a tree of scenes.[^nodes] This skill covers the node tree, the **lifecycle
callback order**, accessing nodes, instancing scenes at runtime, and freeing them.

> Scope: structure and lifecycle only. For events use the `signals` skill; this skill does not cover physics,
> UI or persistence.

## The scene tree & the root

At runtime every node lives under the **SceneTree**. `get_tree().root` is the tree's root node, a
`Window`[^scenetree] (a `Viewport` subclass); your scenes are added beneath it. The **main scene** (set in
Project Settings, stored in `project.godot`) is what loads first.[^nodes] A scene always has exactly one root
node.[^nodes]

## Node lifecycle & callback ordering

This is the part that bites people. When a node enters the tree, callbacks fire in this documented order:[^node]

> `_enter_tree()` of the parent → `_enter_tree()` of the children → `_ready()` of the children → `_ready()`
> of the parent.

So:
- **`_enter_tree()`** runs **top-down** (parent before children).
- **`_ready()`** runs **bottom-up / post-order** — children are ready before their parent, once per node.[^node]
- **`_exit_tree()`** runs **bottom-up**: a parent's `_exit_tree()` is called after all its children have left.[^node]
- **`_process(delta)`** / **`_physics_process(delta)`** run per rendered frame / per physics tick.[^node][^scenetree2]

**Consequence:** a child's nodes are guaranteed to exist in that child's `_ready()`, and a parent can safely
touch its children in the parent's `_ready()` (children are ready first). Do **not** rely on siblings or
parents being ready inside `_enter_tree()`.

```gdscript
extends Node
func _enter_tree() -> void:
	# Parent runs this BEFORE its children's _enter_tree(). Children may not exist yet.
	pass
func _ready() -> void:
	# Children are already ready here. Safe to access them.
	pass
```

**C#** uses PascalCase and `double` for delta; there is **no `@onready`** — assign child references in
`_Ready()`:[^csbasics][^csdiff]

```csharp
public override void _Ready()
{
	_sprite = GetNode<Sprite2D>("Sprite2D"); // no @onready in C#
}
public override void _Process(double delta) { }
```

## Accessing nodes

- `get_node(path)` / the `$` shorthand (GDScript): relative (from this node) or absolute (from
  `/root`).[^node][^instances] `get_node_or_null(path)` returns `null` instead of erroring.[^node]
- **Scene-unique names** `%Name`: mark a node "Access as Unique Name", then `%Name` (or `get_node("%Name")`)
  from any node in the **same scene** — robust against moving the node in the tree.[^unique]

```gdscript
@onready var hp_bar := $UI/HealthBar        # path-based (breaks if you move UI/HealthBar)
@onready var hp_bar2 := %HealthBar          # scene-unique name (survives re-parenting in the scene)
var maybe := get_node_or_null("Optional")
```

**C#:** use `GetNode<T>(path)` (and `%Name`); the `$` shorthand is **GDScript-only**.[^instances][^csdiff]

```csharp
_hpBar = GetNode<TextureProgressBar>("%HealthBar");
```

## Building the tree in code

- `add_child(node)` parents a node at runtime; each child must have a unique name within its
  parent.[^node]
- `remove_child(node)` detaches **without** freeing it (and sets the removed node's `owner` to `null`).[^node]
- `owner` matters for **saving/packing** a scene (only owned nodes are packed), **not** for runtime
  parenting — runtime instancing only needs `add_child`.[^node][^packed]

## Instancing scenes (PackedScene)

A `.tscn`/`PackedScene` is a **resource** (a blueprint), not a node. Load it, then `instantiate()` a fresh
node tree and add it:[^instances][^packed]

```gdscript
const Bullet := preload("res://bullet.tscn")   # preload: GDScript-only, resolved at compile time
func shoot() -> void:
	var b := Bullet.instantiate()              # returns a Node tree
	add_child(b)
```

**C#:** `preload` does **not** exist — use `GD.Load<PackedScene>` or `ResourceLoader.Load<PackedScene>`, then
`Instantiate()`/`Instantiate<T>()`:[^csdiff][^packed]

```csharp
private PackedScene _bullet = GD.Load<PackedScene>("res://bullet.tscn");
public void Shoot() => AddChild(_bullet.Instantiate<Node2D>());
```

## Freeing nodes

- **`queue_free()`** queues deletion at the **end of the current frame** — the node is still valid until
  then; freeing a node also frees its children.[^node]
- **`free()`** deletes **immediately** — unsafe to call while the node (or the engine) is mid-callback; the
  docs **recommend `queue_free()`**.[^instances]
- After deletion, references are invalid; check with **`is_instance_valid(node)`** before using a reference
  you are unsure about.[^globalscope]

```gdscript
enemy.queue_free()                 # removed at end of frame
if is_instance_valid(enemy):
	enemy.position += Vector2.RIGHT
```

## Composition (aggregation), briefly

Godot node trees form an **aggregation** relationship, **not** composition — "nodes are not
components".[^org] Design scenes with few dependencies: a parent configures its children (dependency
injection), and siblings talk through a common ancestor rather than reaching across the tree.[^org] Use a
parent-child relationship when removing the parent should remove the children too.[^org]

## GDScript vs C# at a glance

| Topic | GDScript | C# |
|---|---|---|
| Lifecycle | `_ready`, `_enter_tree`, `_process(delta: float)` | `_Ready`, `_EnterTree`, `_Process(double delta)` |
| On-ready refs | `@onready var x := $X` | assign in `_Ready()` (no `@onready`) |
| Get node | `get_node("X")` / `$X` / `%X` | `GetNode<T>("X")` / `%X` (no `$`) |
| Load scene | `preload`/`load` | `GD.Load<PackedScene>` (no `preload`) |
| Instance | `scene.instantiate()` | `scene.Instantiate<T>()` |
| Free | `queue_free()` | `QueueFree()` |

C# also requires the **.NET editor build** and **cannot export to Web**; use **.NET 8** for current Godot
4.7.[^csindex][^csbasics] More in [references/csharp-scenes.md](references/csharp-scenes.md).

## Common mistakes this skill prevents

- Accessing a child inside `_enter_tree()` (it may not exist yet) instead of `_ready()`.
- Relying on `$Fixed/Path` that breaks when a node moves — prefer `%UniqueName`.
- Calling `free()` during a callback/signal instead of `queue_free()`.
- Assuming C# has `@onready`, `$`, or `preload` (it has none).
- Forgetting `owner` when **building a scene in code to save it** (runtime `add_child` does not need it).

See also [references/lifecycle.md](references/lifecycle.md) and [references/instancing.md](references/instancing.md).

## Examples (and how they are verified)

- [examples/gdscript/node_access.gd](examples/gdscript/node_access.gd) — complete (parse).
- [examples/gdscript/build_and_instance.gd](examples/gdscript/build_and_instance.gd) — complete (parse).
- [examples/gdscript-runtime/](examples/gdscript-runtime/project.godot) — executable; asserts `_ready`
  ordering, `add_child`, instancing and `queue_free` headless (runtime).
- [examples/csharp/](examples/csharp/ScenesDemo.cs) — complete project; `dotnet build` only (BUILD_VERIFIED).

## Sources

- Nodes & scenes: <https://docs.godotengine.org/en/4.7/getting_started/step_by_step/nodes_and_scenes.html>
- Instancing: <https://docs.godotengine.org/en/4.7/getting_started/step_by_step/instancing.html>
- Nodes & scene instances (`$`, get_node, load/preload/instantiate, queue_free): <https://docs.godotengine.org/en/4.7/tutorials/scripting/nodes_and_scene_instances.html>
- Scene tree (lifecycle ordering, root): <https://docs.godotengine.org/en/4.7/tutorials/scripting/scene_tree.html>
- Scene-unique nodes (`%Name`): <https://docs.godotengine.org/en/4.7/tutorials/scripting/scene_unique_nodes.html>
- Scene organization (aggregation): <https://docs.godotengine.org/en/4.7/tutorials/best_practices/scene_organization.html>
- `Node` (lifecycle, get_node, add_child, owner, queue_free): <https://docs.godotengine.org/en/4.7/classes/class_node.html>
- `SceneTree` (root is a Window): <https://docs.godotengine.org/en/4.7/classes/class_scenetree.html>
- `PackedScene` (instantiate, pack, owner & saving): <https://docs.godotengine.org/en/4.7/classes/class_packedscene.html>
- `@GlobalScope` (is_instance_valid): <https://docs.godotengine.org/en/4.7/classes/class_@globalscope.html>
- C# differences (no @onready/$/preload; GetNode&lt;T&gt;; GD.Load): <https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html>
- C# basics (.NET build, PascalCase, no Web): <https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_basics.html>
- C# index (platform constraints): <https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/index.html>

[^nodes]: Nodes & scenes — node/scene model, single root, main scene.
[^node]: `Node` class reference — lifecycle callbacks + the exact enter/ready ordering, `get_node`, `add_child`, `owner`, `queue_free`.
[^scenetree]: `SceneTree` class reference — `root` is a `Window`.
[^scenetree2]: Scene tree tutorial — process pre-order, `_ready` post-order.
[^instances]: Nodes & scene instances tutorial — `$`, `get_node`, `load`/`preload`, `instantiate`, prefer `queue_free`.
[^unique]: Scene-unique nodes tutorial — `%Name` syntax and same-scene constraint.
[^packed]: `PackedScene` class reference — `instantiate()`, `pack()`, owner & saving.
[^globalscope]: `@GlobalScope` — `is_instance_valid()`.
[^org]: Scene organization — aggregation (not composition); dependency injection.
[^csdiff]: C# differences — no `@onready`/`$`/`preload`; `GetNode<T>`; `GD.Load<PackedScene>`.
[^csbasics]: C# basics — PascalCase, `_Process(double)`, .NET-enabled editor build, no Web export.
[^csindex]: C# index — platform constraints (no Web; experimental Android/iOS).
