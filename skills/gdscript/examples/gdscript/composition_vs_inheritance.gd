extends RefCounted

## Complete, parse-verifiable: inheritance for a true "is-a", composition for "has-a".

# Inheritance — a Rifle IS A Weapon, so extending is appropriate.
class Weapon:
	var damage: int = 10
	func fire() -> int:
		return damage

class Rifle extends Weapon:
	func _init() -> void:
		damage = 25

# Composition — an Actor HAS A Weapon (not "is-a"). Holding it as a member lets the
# actor swap weapons without changing its type hierarchy.
class Actor:
	var _weapon: Weapon

	func equip(w: Weapon) -> void:
		_weapon = w

	func attack() -> int:
		return _weapon.fire() if _weapon != null else 0
