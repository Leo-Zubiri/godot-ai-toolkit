extends Node

## Complete, parse-verifiable example: declare custom signals, emit them, and
## connect with a bound extra argument (which arrives AFTER the signal's args).

signal health_changed(old_value, new_value)
signal died

var _health: int = 100

func _ready() -> void:
	health_changed.connect(_on_health_changed.bind("player"))
	died.connect(_on_died)

func damage(amount: int) -> void:
	var old := _health
	_health = maxi(_health - amount, 0)
	health_changed.emit(old, _health)
	if _health == 0:
		died.emit()

func _on_health_changed(old_value: int, new_value: int, who: String) -> void:
	print("%s: %d -> %d" % [who, old_value, new_value])

func _on_died() -> void:
	print("died")
