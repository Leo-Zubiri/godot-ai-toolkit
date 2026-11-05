---
name: signals
description: >-
  Godot 4.7 signals for decoupled node communication in GDScript and C#. Use when
  wiring "something happened" events between nodes/scenes — declaring, emitting,
  connecting/disconnecting (code or editor), Callable/bind, one-shot and deferred
  flags, duplicate-connection safety, awaiting signals, and lifecycle/cleanup.
---

# Signals (Godot 4.7)

Signals are Godot's built-in observer mechanism: an object emits a named event and any number of
connected `Callable`s run. They are the idiomatic way to keep nodes **decoupled** — an emitter does not
need to know who listens. In Godot 4 signals are a first-class type passed without strings.[^signals]

> Scope: this skill is only about signals. For scene/node structure see a future `scenes-and-nodes` skill.

## When to use signals (and when not)

- **Use** for one-to-many, "**past-tense**" notifications (`died`, `health_changed`, `pressed`) where the
  emitter should not depend on its listeners.
- **Recommendation (toolkit, not a cited rule):** "**signal up, call down**" — a child notifies its parent
  via a signal; a parent commands a child via a direct method call. Prefer a direct call for imperative
  commands (`kill()`), and a direct reference for data you can just read.
- Signals add indirection; for tight, per-frame, parent→child logic a direct call is clearer.

## Built-in signals

Most nodes expose signals already, e.g. `Button`'s `pressed`,[^basebutton] `Timer`'s `timeout`,[^timer]
and every `Node`'s `ready`.[^node] Connect to them exactly like custom ones.

## Declaring a custom signal

**GDScript** — the `signal` keyword; argument names are optional and shown in the editor's Signals dock:[^gds]

```gdscript
extends Node
signal health_changed(old_value, new_value)
signal died
```

**C#** — the `[Signal]` attribute on a `public delegate` whose name **must end with `EventHandler`**. Godot
generates a C# `event` (name minus `EventHandler`) and a `SignalName.<Name>` entry; **(re)build** the
project for it to appear:[^cssig]

```csharp
public partial class Health : Node
{
	[Signal] public delegate void HealthChangedEventHandler(int current, int maximum);
}
// -> generated: event HealthChanged, SignalName.HealthChanged
```

> Flag: typed GDScript signal parameters (`signal s(x: int)`) are **not** shown in the 4.7 docs, so this
> skill teaches untyped declarations. Typed *callbacks* (the receiving function) **are** documented.[^gds]

## Emitting

```gdscript
health_changed.emit(old_value, new_value)   # idiomatic 4.x
# emit_signal("health_changed", old_value, new_value)  # also valid; returns an Error
```
`Signal.emit()` takes a variable number of args; `Object.emit_signal()` returns an `Error`.[^signalcls][^object]

```csharp
EmitSignal(SignalName.HealthChanged, current, maximum);
```
You **cannot** use C#'s `Invoke` to raise a Godot-bound event.[^cssig]

## Connecting & disconnecting in code

**GDScript** — connect a `Callable` (the first-class form). Bind extra args with `.bind()`; bound values
arrive **after** the signal's own arguments:[^gds][^signalcls][^callable]

```gdscript
button.pressed.connect(_on_pressed)
health.health_changed.connect(_on_health_changed.bind(player.name))  # name arrives last
# Object form + explicit Callable:
enemy.connect("died", Callable(self, "_on_enemy_died"))
button.pressed.disconnect(_on_pressed)
```
The Godot 3 string form `connect("sig", self, "method")` no longer exists — `connect` takes a `Callable`.[^object]

**C#** — the generated event with `+=` / `-=` is recommended; use `Connect(...)` when connecting to a
GDScript-defined signal or passing flags. C# has **no `bind`/`unbind` — use lambdas**:[^cssig][^csdiff]

```csharp
public override void _Ready()
{
	_health.HealthChanged += OnHealthChanged;              // recommended
	button.Connect(Button.SignalName.Pressed,             // older API / flags / GDScript signals
		Callable.From(OnPressed), (uint)GodotObject.ConnectFlags.OneShot);
}
public override void _ExitTree() => _health.HealthChanged -= OnHealthChanged;
```

## Editor connections

You can also connect a signal to a method through the editor's **Node → Signals** dock; editor-made
connections are persisted with the scene (`CONNECT_PERSIST`).[^signals][^object]

## Connect flags (`Object.ConnectFlags`)

| Flag (GDScript / C#) | Value | Effect |
|---|---|---|
| `CONNECT_DEFERRED` / `ConnectFlags.Deferred` | 1 | Callable runs at idle time (end of frame), not instantly. |
| `CONNECT_PERSIST` / `ConnectFlags.Persist` | 2 | Connection is stored on serialization. |
| `CONNECT_ONE_SHOT` / `ConnectFlags.OneShot` | 4 | Disconnects itself after one emission. |
| `CONNECT_REFERENCE_COUNTED` / `ConnectFlags.ReferenceCounted` | 8 | Allows/ref-counts duplicate connects. |

Values and GDScript names are from the class reference; only `OneShot` is shown in C# form in the docs —
the others follow the standard PascalCase mapping.[^object][^cssig]

## Duplicate-connection safety

A signal can be connected to the same `Callable` **only once**; a duplicate connect returns
`ERR_INVALID_PARAMETER` and logs an error — unless `CONNECT_REFERENCE_COUNTED`. `disconnect()` errors if the
connection doesn't exist. Guard with `is_connected()`:[^object][^signalcls]

```gdscript
if not button.pressed.is_connected(_on_pressed):
	button.pressed.connect(_on_pressed)
```

## Awaiting a signal

```gdscript
await $Button.button_up          # single arg -> that value; many -> Array; none -> null
```
A function that `await`s becomes a coroutine; callers must `await` it too.[^gds] In C#:
`await ToSignal(timer, Timer.SignalName.Timeout);`.[^csdiff]

## Lifecycle & cleanup (avoid dangling connections)

Relevant `Node` signals: `tree_entered`, `tree_exiting`, `tree_exited`, `ready`, `renamed`.[^node]

- **GDScript:** it is documented that **if the callable's object is freed, that connection is lost**.[^object]
  *Recommendation (toolkit):* for a long-lived emitter connected to shorter-lived receivers, disconnect in
  `_exit_tree()` (or use `CONNECT_ONE_SHOT`) rather than relying on broad auto-cleanup.
- **C#:** freeing an object auto-disconnects its connections, **but** a **capturing lambda** or a **custom
  signal connected via `+=`** does **not** auto-disconnect when the receiver is freed — accessing it later
  throws `ObjectDisposedException`. Fix: store the delegate and `-=` it in `_ExitTree`, or connect with
  `Connect` (which auto-disconnects for custom signals).[^cssig]

Details in [references/lifecycle-and-safety.md](references/lifecycle-and-safety.md).

## GDScript vs C# at a glance

| Topic | GDScript | C# |
|---|---|---|
| Declare | `signal x(a, b)` | `[Signal] delegate void XEventHandler(...)` (+ rebuild) |
| Emit | `x.emit(a, b)` | `EmitSignal(SignalName.X, a, b)` |
| Connect | `x.connect(cb)` | `X += cb` / `Connect(SignalName.X, Callable.From(cb))` |
| Bind args | `cb.bind(extra)` | lambda (no `bind`) |
| Await | `await obj.x` | `await ToSignal(obj, T.SignalName.X)` |
| Cleanup | freeing callable's object drops it | `+=`/capturing lambda need manual `-=` |

C# also requires the **.NET editor build** and **cannot export to Web**. **Use .NET 8** for current Godot
4.7 (9 for Android export); .NET 6 was the technical minimum only for Godot 4.0–4.3 and is **not recommended**
for new work.[^csindex][^csbasics] More in [references/csharp-signals.md](references/csharp-signals.md).

## Anti-patterns

See [references/anti-patterns.md](references/anti-patterns.md) — e.g. using signals for synchronous
request/response, connecting without a disconnect strategy, or emitting to drive imperative commands.

## Examples (and how they are verified)

- [examples/gdscript/declare_emit_connect.gd](examples/gdscript/declare_emit_connect.gd) — complete (parse).
- [examples/gdscript/safe_connect.gd](examples/gdscript/safe_connect.gd) — complete (parse).
- [examples/gdscript-runtime/](examples/gdscript-runtime/project.godot) — executable; asserts real emit /
  connect / one-shot / duplicate-guard behavior headless (runtime).
- [examples/csharp/](examples/csharp/SignalsDemo.cs) — complete project; `dotnet build` only (BUILD_VERIFIED).

## Sources

- GDScript basics (signal keyword, await): <https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/gdscript_basics.html>
- Signals tutorial: <https://docs.godotengine.org/en/4.7/getting_started/step_by_step/signals.html>
- `Object` (connect/disconnect/emit_signal/is_connected, ConnectFlags): <https://docs.godotengine.org/en/4.7/classes/class_object.html>
- `Signal`: <https://docs.godotengine.org/en/4.7/classes/class_signal.html>
- `Callable`: <https://docs.godotengine.org/en/4.7/classes/class_callable.html>
- `Node` (lifecycle signals): <https://docs.godotengine.org/en/4.7/classes/class_node.html>
- `BaseButton` / `Timer` (built-in signals): <https://docs.godotengine.org/en/4.7/classes/class_basebutton.html>, <https://docs.godotengine.org/en/4.7/classes/class_timer.html>
- C# signals: <https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_signals.html>
- C# differences: <https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html>
- C# basics (.NET build/SDK, no Web): <https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_basics.html>
- C# index (platform constraints): <https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/index.html>

[^signals]: Signals tutorial — first-class signals since Godot 4.0.
[^gds]: GDScript basics — `signal`, `.emit()`, `.connect()`, `.bind()`, `await`.
[^signalcls]: `Signal` class reference — `emit`, `connect`, `disconnect`, `is_connected`.
[^object]: `Object` class reference — `connect`/`disconnect`/`emit_signal`/`is_connected`, `ConnectFlags`.
[^callable]: `Callable` class reference — `Callable(obj, method)`, `bind`.
[^node]: `Node` class reference — `ready`, `tree_entered`, `tree_exiting`, `tree_exited`, `renamed`.
[^basebutton]: `BaseButton` class reference — `pressed`, `button_up`, `button_down`, `toggled`.
[^timer]: `Timer` class reference — `timeout`.
[^cssig]: C# signals — `[Signal]` + `EventHandler`, `EmitSignal(SignalName…)`, `+=`/`-=`, `Connect`, lifecycle.
[^csdiff]: C# differences — no `bind`/`unbind` (use lambdas), `ToSignal`, `Singleton`.
[^csindex]: C# index — C# cannot export to Web; experimental Android/iOS.
[^csbasics]: C# basics — .NET-enabled editor build, .NET 8 SDK floor (9 for Android).
