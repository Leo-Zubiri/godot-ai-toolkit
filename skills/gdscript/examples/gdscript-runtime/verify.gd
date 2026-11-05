extends SceneTree

## Headless runtime verification of real GDScript-language behavior (4.0-safe):
##  - typed Array[int] holds and returns ints
##  - Dictionary is passed BY REFERENCE (a function mutates the caller's value)
##  - named enum values
##  - is / as type checks
##  - inner class method
##  - lambda captures locals by value at creation
## Prints "FAIL: ..." and quits 1 on mismatch; "OK" and quits 0 on success.
## Run: godot --headless --path . --script res://verify.gd

enum Dir { LEFT, RIGHT, UP, DOWN }

class Counter:
	var value: int = 0
	func inc() -> void:
		value += 1

func _initialize() -> void:
	var failures := 0

	# 1. typed array
	var a: Array[int] = []
	a.append(1)
	a.append(2)
	if a.size() != 2 or a[0] != 1:
		failures += _fail("typed array append; got %s" % [a])

	# 2. dictionary passed by reference
	var d := {"hp": 10}
	_mutate(d)
	if int(d.get("hp")) != 20:
		failures += _fail("Dictionary not passed by reference; hp=%s" % [d.get("hp")])

	# 3. named enum values (LEFT=0, RIGHT=1, ...)
	if Dir.RIGHT != 1:
		failures += _fail("enum Dir.RIGHT=%s, expected 1" % [Dir.RIGHT])

	# 4. is / as
	var n := Node.new()
	if not (n is Node):
		failures += _fail("'is Node' was false")
	if (n as Node2D) != null:
		failures += _fail("'as Node2D' on a plain Node should be null")
	n.free()

	# 5. inner class method
	var c := Counter.new()
	c.inc()
	c.inc()
	if c.value != 2:
		failures += _fail("inner class method; value=%s" % [c.value])

	# 6. lambda captures by value at creation (later reassignment not visible)
	var base := 10
	var add := func(x): return x + base
	base = 999
	if add.call(5) != 15:
		failures += _fail("lambda capture-by-value; got %s" % [add.call(5)])

	if failures == 0:
		print("OK: gdscript runtime checks passed")
		quit(0)
	else:
		quit(1)

func _mutate(dict: Dictionary) -> void:
	dict["hp"] = 20

func _fail(msg: String) -> int:
	print("FAIL: " + msg)
	return 1
