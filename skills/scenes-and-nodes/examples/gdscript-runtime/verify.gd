extends SceneTree

## Headless runtime verification of real scene/node behavior:
##  - _ready callback ordering (children before parent / post-order)
##  - add_child + get_node_or_null (by node name)
##  - PackedScene.pack + instantiate reproduces the owned subtree
##  - free() invalidates the instance (is_instance_valid)
## The tree is built in _initialize(); assertions run in _process() so that the
## runtime-added nodes' _ready() callbacks (which fire on the frame, not inside
## _initialize) have already run. Prints "FAIL: ..." and quits 1 on mismatch;
## "OK" and quits 0 on success.
## Run: godot --headless --path . --script res://verify.gd

var order: Array = []
var parent: Node
var child: Node
var _ran := false

class OrderNode extends Node:
	var tag: String
	var sink
	func _ready() -> void:
		sink.order.append(tag)

func _initialize() -> void:
	parent = OrderNode.new()
	parent.name = "parent"
	parent.tag = "parent"
	parent.sink = self
	child = OrderNode.new()
	child.name = "child"
	child.tag = "child"
	child.sink = self
	parent.add_child(child)
	root.add_child(parent)  # _enter_tree/_ready fire on the frame, not synchronously here

func _process(_delta: float) -> bool:
	if _ran:
		return false
	_ran = true
	var failures := 0

	# 1. _ready ordering: child before parent (post-order).
	if order != ["child", "parent"]:
		failures += _fail("_ready order was %s, expected [child, parent]" % [order])

	# 2. get_node_or_null finds the child by name.
	if parent.get_node_or_null("child") != child:
		failures += _fail("get_node_or_null did not find the child")

	# 3. PackedScene.pack (owned nodes only) + instantiate reproduces the subtree.
	var proot := Node.new()
	proot.name = "proot"
	var leaf := Node.new()
	leaf.name = "leaf"
	proot.add_child(leaf)
	leaf.owner = proot  # required so pack() saves it
	var ps := PackedScene.new()
	ps.pack(proot)
	var inst := ps.instantiate()
	if inst.get_node_or_null("leaf") == null:
		failures += _fail("instantiated scene is missing 'leaf'")
	inst.free()
	proot.free()

	# 4. free() invalidates the reference immediately.
	var tmp := Node.new()
	root.add_child(tmp)
	tmp.free()
	if is_instance_valid(tmp):
		failures += _fail("node still valid after free()")

	parent.free()  # clean up the built tree (no leaked nodes at exit)

	if failures == 0:
		print("OK: scenes-and-nodes runtime checks passed")
		quit(0)
	else:
		quit(1)
	return false

func _fail(msg: String) -> int:
	print("FAIL: " + msg)
	return 1
