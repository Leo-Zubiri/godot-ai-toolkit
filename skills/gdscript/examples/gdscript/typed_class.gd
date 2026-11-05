class_name HealthComponent extends RefCounted

## Complete, parse-verifiable: a typed class with an enum, typed fields and methods
## (4.0-safe — no typed dictionaries, @abstract or static vars).

enum State { ALIVE, DEAD }

var _current: int
var _max: int
var state: State = State.ALIVE

func _init(maximum: int = 100) -> void:
	_max = maximum
	_current = maximum

func damage(amount: int) -> void:
	_current = maxi(_current - amount, 0)
	if _current == 0:
		state = State.DEAD

func ratio() -> float:
	return float(_current) / float(_max)
