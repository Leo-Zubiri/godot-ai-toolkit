# C# signals — details

Companion to `SKILL.md`. All claims cited to the Godot 4.7 C# docs (and Microsoft for the underlying
event/delegate model, labeled as secondary support).

## The generated-event model

A `[Signal]` delegate named `<Name>EventHandler` makes Godot generate a C# `event` named `<Name>` plus a
`SignalName.<Name>` entry. The project must be **(re)built** for new signals to appear in the editor.
Godot's binding sits on top of the standard C# **delegate → event** model (publisher/subscriber).
Sources: [C# signals](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_signals.html);
secondary: [MS events](https://learn.microsoft.com/en-us/dotnet/csharp/programming-guide/events/),
[MS delegates](https://learn.microsoft.com/en-us/dotnet/csharp/programming-guide/delegates/).

```csharp
public partial class Health : Node
{
	[Signal] public delegate void HealthChangedEventHandler(int current, int maximum);

	private int _current = 100;

	public void Damage(int amount)
	{
		_current -= amount;
		EmitSignal(SignalName.HealthChanged, _current, 100); // not Invoke()
	}
}
```

## Connecting: `+=` vs `Connect`

- `HealthChanged += OnChanged;` / `-=` — **recommended** for your own signals.
- `node.Connect(SignalName.X, Callable.From(OnX), (uint)GodotObject.ConnectFlags.OneShot);` — use when
  connecting to a **GDScript-defined** signal or when you need `ConnectFlags`.
- **No `bind`/`unbind` in C#** — capture with a lambda: `Callable.From(() => SayHello(name))`.

Source: [C# signals](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_signals.html),
[C# differences](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html).

## Awaiting

```csharp
await ToSignal(timer, Timer.SignalName.Timeout);
```
Source: [C# differences](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_differences.html).

## Naming & identifiers

The C# API is PascalCase. String-based calls (`Connect`, `Call`, …) still expect Godot's snake_case names —
prefer the generated `SignalName`/`MethodName` identifiers to avoid typos and `StringName` allocations.
Source: [C# basics](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_basics.html).

## Platform / toolchain constraints (verified)

- C# requires the **.NET-enabled editor build** and the **.NET SDK**.
- **.NET 8** is the SDK floor; **Android export requires .NET 9**.
- **C# cannot export to the Web platform.**

Sources: [C# index](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/index.html),
[C# basics](https://docs.godotengine.org/en/4.7/tutorials/scripting/c_sharp/c_sharp_basics.html).

> Flag: the 4.7 `c_sharp_basics.html` text literally reads "Godot 4.5 requires .NET 8…" — a stale version
> label. The **requirement** (.NET 8, 9 for Android) is what the page states; we do not quote it as
> "4.7 requires".
