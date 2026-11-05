# C# scenes & nodes — details

Companion to `SKILL.md`. Cited to the Godot 4.7 C# docs.

## Lifecycle

`_Ready()` and `_Process(double delta)` are documented in
[c_sharp_basics.html](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_basics.html).
There is **no `@onready`** — assign node references in `_Ready()`
([c_sharp_differences.html](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html)).
`_EnterTree()`/`_ExitTree()`/`_PhysicsProcess(double)` are the conventional PascalCase overrides (not
individually quoted on these pages).

## Accessing nodes

- `GetNode<T>("Name")` is the generic form; nodes are fetched by **name/path, not type**.
- The `$` shorthand is **GDScript-only** (no C# equivalent).
- Scene-unique names work as `%Name`.
Source: [nodes & scene instances](https://docs.godotengine.org/en/4.7/tutorials/scripting/nodes_and_scene_instances.html),
[c_sharp_differences.html](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html).

## Instancing & freeing

- `GD.Load<PackedScene>("res://x.tscn")` or `ResourceLoader.Load<PackedScene>(...)`; `preload` is unavailable.
- `Instantiate()` / `Instantiate<T>()`, then `AddChild(...)`.
- `QueueFree()`; check validity with `GodotObject.IsInstanceValid(obj)`.
Source: [c_sharp_differences.html](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html),
[class_packedscene.html](https://docs.godotengine.org/en/4.7/classes/class_packedscene.html).

## Naming & constraints

- **PascalCase** for the API; but string-based calls (`Call`, `Connect`, `CallDeferred`, `Get`, `Set`) still
  take Godot's snake_case names (e.g. `CallDeferred("add_child")`).
- C# requires the **.NET-enabled editor build** and the .NET SDK; **C# cannot export to the Web**.
Source: [c_sharp_basics.html](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_basics.html),
[index.html](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/index.html).

> Honesty flags: `GetNodeOrNull<T>` and the implicit `string`→`NodePath` conversion are used in practice but
> are not individually quoted on the 4.7 pages reviewed. Freed-node access throwing `ObjectDisposedException`
> is **not** documented on these pages — use `IsInstanceValid` as the documented safety check.
