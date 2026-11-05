extends Node

## Complete, parse-verifiable example: lifecycle order and node access.

func _enter_tree() -> void:
	# Parent runs this before its children's _enter_tree(); children may not exist yet.
	pass

func _ready() -> void:
	# Children are ready here (post-order). Safe to access them.
	var child := get_node_or_null("Child")       # null-safe lookup
	if child != null:
		print(child.name)
	var hp := get_node_or_null("%HealthBar")      # scene-unique name, survives re-parenting
	if hp != null:
		print(hp)
