extends RefCounted

## Complete, parse-verifiable: typed arrays and a small data model (4.0-safe;
## uses a plain Dictionary since typed dictionaries are Godot 4.4+).

class Item:
	var name: String
	var value: int
	func _init(n: String, v: int) -> void:
		name = n
		value = v

var _items: Array[Item] = []

func add(item: Item) -> void:
	_items.append(item)

func total_value() -> int:
	var sum := 0
	for it in _items:
		sum += it.value
	return sum

func names() -> Array[String]:
	var out: Array[String] = []
	for it in _items:
		out.append(it.name)
	return out

func index_by_name() -> Dictionary:
	var by_name := {}
	for it in _items:
		by_name[it.name] = it
	return by_name
