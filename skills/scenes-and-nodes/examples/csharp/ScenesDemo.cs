using Godot;

// Complete C# example (BUILD_VERIFIED via Godot.NET.Sdk; not executed here).
// Node access, instancing and freeing in C#: no @onready, no $, no preload.
public partial class ScenesDemo : Node
{
	// C# uses GD.Load / ResourceLoader.Load (preload does not exist in C#).
	private PackedScene _bullet = GD.Load<PackedScene>("res://bullet.tscn");
	private Node2D _holder = null!;

	public override void _Ready()
	{
		// Assign node references here — C# has no @onready.
		_holder = GetNodeOrNull<Node2D>("Holder") ?? new Node2D { Name = "Holder" };
		if (_holder.GetParent() == null)
		{
			AddChild(_holder);
		}
	}

	public void Shoot()
	{
		Node2D bullet = _bullet.Instantiate<Node2D>();
		_holder.AddChild(bullet); // runtime parenting needs no owner
	}

	public void Despawn(Node node)
	{
		if (GodotObject.IsInstanceValid(node))
		{
			node.QueueFree();
		}
	}
}
