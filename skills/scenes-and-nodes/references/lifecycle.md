# Node lifecycle — details

Companion to `SKILL.md`. Cited to the Godot 4.7 class reference and scene-tree tutorial.

## The exact ordering (quoted)

From [`class_node.html`](https://docs.godotengine.org/en/4.7/classes/class_node.html): when a node is added to
the tree the callbacks fire as `_enter_tree()` of the parent → `_enter_tree()` of the children → `_ready()` of
the children → `_ready()` of the parent. `_ready()` runs once per node, after all its children are ready;
`_exit_tree()` runs bottom-up (a parent after its children).

The [scene-tree tutorial](https://docs.godotengine.org/en/4.7/tutorials/scripting/scene_tree.html) summarizes
it as: most operations run **pre-order** (top to bottom, e.g. `_process`); `_ready()` runs **post-order**
(children first).

## Notification constants (class_node.html)

- `NOTIFICATION_ENTER_TREE = 10`, `NOTIFICATION_READY = 13`, `NOTIFICATION_EXIT_TREE = 11`,
  `NOTIFICATION_PROCESS = 17`, `NOTIFICATION_PHYSICS_PROCESS = 16`.

## Practical rules

- Access a node's own children in **its** `_ready()` (they are ready first). A parent may also touch its
  children in the parent's `_ready()`.
- Do **not** assume siblings/parents are ready in `_enter_tree()`.
- `_process(delta)` is per rendered frame; `_physics_process(delta)` is per physics tick. `delta` is `float`
  in GDScript and `double` in C#.

## C# note (honesty flag)

The 4.7 C# pages quote `_Ready()` and `_Process(double delta)` directly. The C# forms `_EnterTree()`,
`_ExitTree()` and `_PhysicsProcess(double delta)` follow the documented PascalCase mapping of the GDScript
virtuals (they are the standard C# overrides), but the C# signatures are not individually quoted on the pages
reviewed — treat them as the conventional PascalCase equivalents.
