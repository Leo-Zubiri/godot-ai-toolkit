extends Node

## Complete, parse-verifiable example: build the tree in code, instance a scene,
## and free nodes safely.

func spawn(scene: PackedScene) -> Node:
	var inst := scene.instantiate()   # a fresh node tree (the scene is a resource)
	add_child(inst)                   # runtime parenting needs no owner
	return inst

func despawn(node: Node) -> void:
	if is_instance_valid(node):
		node.queue_free()             # removed at end of frame; prefer over free()

func build_manually() -> void:
	var holder := Node2D.new()
	holder.name = "Holder"
	add_child(holder)
	var sprite := Sprite2D.new()
	holder.add_child(sprite)
