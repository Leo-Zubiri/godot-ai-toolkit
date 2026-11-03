using Godot;

// Minimal, deterministic C# node for the fixture. Demonstrates the Godot C# API
// shape (PascalCase, _Ready override) without any real game logic.
public partial class Hello : Node
{
	public override void _Ready()
	{
		GD.Print("hello from the csharp fixture");
	}
}
