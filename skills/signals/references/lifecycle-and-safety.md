# Signal lifecycle & safety

Companion to `SKILL.md`. Cited to the Godot 4.7 class reference and the C# signals page.

## Node lifecycle signals

`Node` emits `tree_entered`, `ready`, `renamed`, `tree_exiting` (about to leave) and `tree_exited` (left).
Connect/disconnect around these when a connection's validity depends on tree membership.
Source: [`Node`](https://docs.godotengine.org/en/4.7/classes/class_node.html).

## Duplicate connections

Connecting the same `Callable` to a signal twice returns `ERR_INVALID_PARAMETER` and logs an error, unless
you pass `CONNECT_REFERENCE_COUNTED`. `disconnect()` errors if the connection does not exist. Always guard:

```gdscript
if not emitter.changed.is_connected(_on_changed):
	emitter.changed.connect(_on_changed)
```
Source: [`Object`](https://docs.godotengine.org/en/4.7/classes/class_object.html),
[`Signal`](https://docs.godotengine.org/en/4.7/classes/class_signal.html).

## What happens when objects are freed

- **Documented (GDScript/engine):** if the **callable's object is freed**, that connection is lost.
  Source: [`Object.connect`](https://docs.godotengine.org/en/4.7/classes/class_object.html).
- **Toolkit recommendation (not a blanket doc guarantee):** do not rely on implicit cleanup for a
  long-lived emitter. Disconnect in `_exit_tree()`, or use `CONNECT_ONE_SHOT` for fire-once wiring.

## C#-specific hazards (documented)

From [C# signals](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_signals.html):

- Freeing an object auto-disconnects its connections — reliable **if the emitter is freed before its
  receivers**.
- A **capturing lambda** does **not** auto-disconnect; after the receiver is freed, invocation throws
  `System.ObjectDisposedException`. Store the delegate and remove it:

```csharp
private Action _tick;
public override void _EnterTree()  { _tick = () => DoThing(); _timer.Timeout += _tick; }
public override void _ExitTree()   { _timer.Timeout -= _tick; }
```

- A **custom signal connected with `+=`** also does **not** auto-disconnect on receiver free; either `-=` in
  `_ExitTree`, or use `Connect` (which auto-disconnects for custom signals). `Callable.From` does not change
  `Delegate.Target`, so it alone does not fix the capturing-lambda case.
