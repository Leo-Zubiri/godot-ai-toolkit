# Instancing & the node tree — details

Companion to `SKILL.md`. Cited to the Godot 4.7 class reference and tutorials.

## PackedScene is a resource

A `.tscn`/`PackedScene` is a resource (a blueprint), not a node.
[`PackedScene.instantiate()`](https://docs.godotengine.org/en/4.7/classes/class_packedscene.html) "instantiates
the scene's node hierarchy" and returns a `Node` tree; `can_instantiate()` returns `true` if the scene has
nodes. The returned tree is not in the SceneTree until you `add_child()` it.

## Loading

- GDScript: `preload("res://x.tscn")` (resolved at compile time, GDScript-only) or `load("res://x.tscn")`
  (runtime). — [nodes & scene instances](https://docs.godotengine.org/en/4.7/tutorials/scripting/nodes_and_scene_instances.html)
- C#: `preload` is not available — use `GD.Load<PackedScene>("res://x.tscn")` or
  `ResourceLoader.Load<PackedScene>(...)`. — [C# differences](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html)

## Building a tree in code

- `add_child(node)` parents at runtime; each child needs a unique name within its parent.
- `remove_child(node)` detaches without freeing and sets the removed node's `owner` to `null`.
- `owner` is required for **saving/packing**: `PackedScene.pack(root)` saves only the nodes that are *owned*
  (an ancestor is set as their `owner`). For plain runtime instancing you only need `add_child` — no `owner`.
  — [`class_node.html`](https://docs.godotengine.org/en/4.7/classes/class_node.html),
  [`class_packedscene.html`](https://docs.godotengine.org/en/4.7/classes/class_packedscene.html)

## Freeing

- `queue_free()` defers deletion to the end of the frame; the node stays valid until then; freeing frees
  children. `free()` is immediate and the docs recommend `queue_free()` instead.
- After deletion, references are invalid; use
  [`is_instance_valid(node)`](https://docs.godotengine.org/en/4.7/classes/class_@globalscope.html) to check.
  — [`class_node.html`](https://docs.godotengine.org/en/4.7/classes/class_node.html)
