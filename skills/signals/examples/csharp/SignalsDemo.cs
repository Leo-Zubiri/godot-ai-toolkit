using Godot;

// Complete C# example (BUILD_VERIFIED via Godot.NET.Sdk; not executed here).
// Shows: [Signal] declaration, EmitSignal(SignalName...), the recommended
// += / -= event, and Connect(...) with Callable.From + a one-shot flag.
public partial class SignalsDemo : Node
{
	[Signal]
	public delegate void HealthChangedEventHandler(int current, int maximum);

	private int _health = 100;

	public override void _Ready()
	{
		// Recommended: the generated C# event.
		HealthChanged += OnHealthChanged;

		// Connect form — required for GDScript-defined signals or ConnectFlags.
		var timer = new Timer();
		AddChild(timer);
		timer.Connect(
			Timer.SignalName.Timeout,
			Callable.From(OnTimeout),
			(uint)GodotObject.ConnectFlags.OneShot);
	}

	public void Damage(int amount)
	{
		_health = Mathf.Max(_health - amount, 0);
		EmitSignal(SignalName.HealthChanged, _health, 100);
	}

	private void OnHealthChanged(int current, int maximum)
	{
		GD.Print($"health: {current}/{maximum}");
	}

	private void OnTimeout()
	{
		GD.Print("timeout");
	}

	public override void _ExitTree()
	{
		// `+=` does not auto-disconnect when this node is freed — remove it.
		HealthChanged -= OnHealthChanged;
	}
}
