extends SceneTree

## Headless runtime verification of real signal behavior.
## Prints "FAIL: ..." and quits with code 1 on any mismatch; prints "OK" and
## quits 0 on success. Run: godot --headless --path . --script res://verify.gd

signal value_changed(v)

var _received: Array = []
var _oneshot_count: int = 0

func _initialize() -> void:
	var failures := 0

	# 1. Connect + emit delivers the argument.
	value_changed.connect(_on_value_changed)
	value_changed.emit(42)
	if _received != [42]:
		failures += _fail("emit did not deliver arg; got %s" % [_received])

	# 2. is_connected reflects the connection.
	if not value_changed.is_connected(_on_value_changed):
		failures += _fail("is_connected false after connect")

	# 3. One-shot connection fires exactly once, then auto-disconnects.
	value_changed.connect(_on_oneshot, CONNECT_ONE_SHOT)
	value_changed.emit(1)
	value_changed.emit(2)
	if _oneshot_count != 1:
		failures += _fail("one-shot fired %d time(s), expected 1" % _oneshot_count)

	# 4. disconnect stops delivery.
	value_changed.disconnect(_on_value_changed)
	if value_changed.is_connected(_on_value_changed):
		failures += _fail("still connected after disconnect")

	if failures == 0:
		print("OK: signal runtime checks passed")
		quit(0)
	else:
		quit(1)

func _on_value_changed(v) -> void:
	_received.append(v)

func _on_oneshot(_v) -> void:
	_oneshot_count += 1

func _fail(msg: String) -> int:
	print("FAIL: " + msg)
	return 1
