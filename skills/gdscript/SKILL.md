---
name: gdscript
description: >-
  Writing correct, idiomatic, maintainable GDScript for Godot 4.x — static typing
  and inference, classes/inheritance/inner classes, enums, typed collections and
  data modeling, annotations, and defensive programming. Use for language-level
  decisions; for signals or node/scene APIs use those skills instead.
---

# GDScript (the language)

This skill is about **the GDScript language itself** — how to write typed, idiomatic, maintainable code and
choose the right construct. It is **not** about signals or nodes/scenes:

- Signals (`signal`/`emit`/`connect`) → the **`signals`** skill.
- Nodes, scenes, lifecycle, `@onready`, `get_node`, instancing → the **`scenes-and-nodes`** skill.

Versions: examples are verified on **Godot 4.0** (minimum) and **4.7.2** (reference). Features newer than 4.0
are **labeled with their version** below; do not use them if you must support 4.0.

## Static typing & inference

Prefer static types — they catch errors before running and let GDScript emit **optimized opcodes** when types
are known at compile time.[^typing] Declare a type with a colon, infer it with `:=`, and type parameters and
return values:[^typing]

```gdscript
var health: int = 100          # explicit
var speed := 220.0             # inferred (float) — use when the type is obvious on the line
func damage(amount: int) -> void:
	health -= amount
func is_alive() -> bool:
	return health > 0
```

- `Variant` is "not an actual type" — it means *no* static type; inference from a `Variant` is disallowed by
  default.[^basics]
- **Recommendation (style guide):** prefer `:=` when the type is written on the same line as the assignment,
  otherwise write the type explicitly; some values (e.g. `get_node()` results) can't be inferred.[^style]

## Type safety: `is`, `as`, `typeof`

```gdscript
if node is Node2D:                 # true runtime type check
	node.position += Vector2.RIGHT
var n2d := node as Node2D          # cast; null on mismatch for Objects (errors for built-in types)
```

`is` / `is not` test the runtime type; `as` casts and yields `null` on an Object type mismatch (but **errors**
when casting between incompatible *built-in* types).[^typing] `typeof(value)` returns the Variant type id.[^typeof]
Static typing lets Godot "detect more errors without even running the code."[^typing]

## Functions & scope

- Only **trailing** parameters may have defaults.[^basics]
- A non-`void` function must return on every path; a `void` function may `return` early but not a value.[^basics]
- `static func` has no `self`/instance access; **`static var`** is shared across instances — *Godot 4.1+*.[^basics]
- Scope resolves local → class member → global.[^basics]

```gdscript
func give(item: String, qty: int = 1) -> void:
	for i in qty:
		_inventory.append(item)
```

## Classes, inheritance & inner classes

`class_name` registers a global class; omitting `extends` makes the base `RefCounted`. Inner classes use
`class Name:` and instantiate with `.new()`. The constructor is `_init`; call the parent with
`super(...)` (or `super.method()` for a specific method). If a base `_init` takes arguments, the child **must**
also define `_init`.[^basics]

```gdscript
class_name Enemy extends Node

class Weapon:                      # inner class
	var damage: int
	func _init(d: int) -> void:
		damage = d

func _init() -> void:
	_weapon = Weapon.new(10)
```

- **No multiple inheritance.**[^basics]
- **Abstract classes/methods** use `@abstract` — *Godot 4.5+*; an abstract class can't be instantiated.[^basics]
- **Composition vs inheritance (recommendation):** prefer composition — extend a base type only for a true
  "is-a" relationship; otherwise hold collaborators as typed members and delegate. (Scene-level composition is
  covered by `scenes-and-nodes`.)

## Constants & enums

`const` needs a compile-time value. Use enums instead of magic numbers; a **named** enum is a constant
Dictionary accessed as `Name.MEMBER`:[^basics]

```gdscript
const MAX_SPEED := 400.0
enum State { IDLE, RUN, JUMP = 10, FALL }   # IDLE=0, RUN=1, JUMP=10, FALL=11
```

## Built-in types & Variant — value vs reference (important)

Base types (`int`, `float`, `String`, vectors, …) are **copied** on assignment and when passed to functions.
**`Object`, `Array`, `Dictionary` and packed arrays are passed by reference** — mutating them inside a function
mutates the caller's value.[^basics][^array][^dict]

```gdscript
func _reset(stats: Dictionary) -> void:
	stats["hp"] = 100              # caller sees this change (passed by reference)
```

`int`/`float` are 64-bit; only `Object`-derived values can be `null`.[^basics]

## Collections & data modeling

Type your collections for safety and clarity:[^typing][^array]

```gdscript
var scores: Array[int] = [10, 20, 30]     # element type checked on writes, [], and +
var names: Array[String] = []
```

- Nested typed arrays (`Array[Array[int]]`) are **not** supported.[^typing]
- **Typed dictionaries** `Dictionary[String, int]` — *Godot 4.4+*; omit the types (plain `Dictionary`) on 4.0.[^basics]
- Useful methods: `map`/`filter`/`find`/`has`/`sort`/`duplicate(true)` on arrays; `get(key, default)`/`has`/
  `keys`/`merge` on dictionaries.[^array][^dict]

Illustrative (4.4+):

```gdscript
var costs: Dictionary[String, int] = { "sword": 50, "shield": 30 }   # typed dictionary — Godot 4.4+
```

## Annotations (language-level)

`@export` (and variants / `@export_group`) exposes a property in the Inspector; `@tool` runs the script in the
editor; `@static_unload` resets static vars (*4.1+*); `@warning_ignore("name")` silences one warning.[^basics][^warnings]
(`@onready` → `scenes-and-nodes`; `@rpc` → multiplayer.)

## Error handling & defensive programming

GDScript has **no exceptions / try-catch**. Guard with types and checks, and report problems with
`push_error()` / `push_warning()`:[^globalscope]

```gdscript
func equip(item) -> bool:
	if item == null or not (item is Weapon):
		push_error("equip() expected a Weapon")
		return false
	_weapon = item
	return true
```

`assert(condition, message)` checks invariants **only in debug builds — it is ignored in release**, so never
put side effects inside it.[^basics]

```gdscript
assert(health >= 0, "health went negative")   # debug-only
```

## Idioms & style

Files/members `snake_case`, classes `PascalCase`, constants and enum members `CONSTANT_CASE`; use `and`/`or`/
`not` (not `&&`/`||`/`!`); tabs for indentation.[^style]

## Choosing between constructs (quick guide)

| Decision | Prefer | When |
|---|---|---|
| Typed vs untyped | **typed** (`:=` or explicit) | almost always — safety + speed; untyped only for genuinely dynamic data |
| Composition vs inheritance | **composition** | unless it's a true "is-a"; avoid deep hierarchies |
| Array vs Dictionary | **Array** for ordered/indexed; **Dictionary** for keyed lookup | model by access pattern |
| enum vs const | **enum** for a set of related states; **const** for single values | avoid magic numbers |
| lambda vs method | **method** for reused/named logic; **lambda** for short local callbacks | lambdas capture locals by value |

## Common anti-patterns

- Putting side effects in `assert(...)` (they vanish in release builds).
- Assuming `Array`/`Dictionary` are copied when passed — they are **by reference**.
- Using version-gated syntax (typed dicts, `@abstract`) while claiming 4.0 support.
- Leaving collections untyped where a type would catch bugs.
- Reaching for exceptions/try-catch (GDScript has none) instead of return values + `is`/`as` guards.

## Examples (and how they are verified)

- [examples/gdscript/typed_class.gd](examples/gdscript/typed_class.gd) — complete (parse).
- [examples/gdscript/data_modeling.gd](examples/gdscript/data_modeling.gd) — complete (parse).
- [examples/gdscript/composition_vs_inheritance.gd](examples/gdscript/composition_vs_inheritance.gd) — complete (parse).
- [examples/gdscript/defensive.gd](examples/gdscript/defensive.gd) — complete (parse).
- [examples/gdscript-runtime/](examples/gdscript-runtime/project.godot) — executable; asserts typed arrays,
  dictionary pass-by-reference, enums, `is`/`as`, an inner class, and lambda capture headless (runtime).

See also [references/typing-and-classes.md](references/typing-and-classes.md) and
[references/collections-and-errors.md](references/collections-and-errors.md).

## Sources

- GDScript basics (vars/const/func/classes/enums/Variant/arrays/dicts/match/lambdas/annotations/assert):
  <https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/gdscript_basics.html>
- Static typing (typed vars/inference/params/returns, typed arrays, `is`/`as`, performance):
  <https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/static_typing.html>
- Style guide (naming, `:=` guidance, boolean operators): <https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/gdscript_styleguide.html>
- Warning system (`@warning_ignore`): <https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/warning_system.html>
- `@GDScript` (annotations, built-ins, `@abstract`): <https://docs.godotengine.org/en/4.7/classes/class_@gdscript.html>
- `@GlobalScope` (`push_error`/`push_warning`/`typeof`): <https://docs.godotengine.org/en/4.7/classes/class_@globalscope.html>
- `Array` (passed by reference; methods): <https://docs.godotengine.org/en/4.7/classes/class_array.html>
- `Dictionary` (passed by reference; methods): <https://docs.godotengine.org/en/4.7/classes/class_dictionary.html>

[^typing]: Static typing — typed declarations/inference/params/returns, `is`/`as`, typed arrays, performance (optimized opcodes) in 4.7.
[^basics]: GDScript basics — functions/scope, `static var` (4.1+), classes/`_init`/`super`/inner classes, `@abstract` (4.5+), enums, Variant, value-vs-reference, typed dictionaries (4.4+), `assert` (debug-only).
[^style]: GDScript style guide — naming conventions, `:=` guidance, `and`/`or`/`not`.
[^warnings]: Warning system — `@warning_ignore` and promoting warnings to errors.
[^typeof]: `@GlobalScope` — `typeof(variable: Variant) -> int`.
[^globalscope]: `@GlobalScope` — `push_error(...)`, `push_warning(...)`.
[^array]: `Array` class reference — "always passed by reference"; `map`/`filter`/`find`/`sort`/`duplicate`.
[^dict]: `Dictionary` class reference — "always passed by reference"; `get`/`has`/`keys`/`merge`.
