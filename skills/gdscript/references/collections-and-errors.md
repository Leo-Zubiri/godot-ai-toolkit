# Collections, value/reference & error handling — details

Companion to `SKILL.md`. Cited to the Godot 4.7 docs.

## Value vs reference

Base types (`int`, `float`, `String`, vectors, …) are copied on assignment and when passed; `Object`, `Array`,
`Dictionary` and packed arrays are **passed by reference**.
[basics](https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/gdscript_basics.html),
[`Array`](https://docs.godotengine.org/en/4.7/classes/class_array.html) ("always passed by reference"),
[`Dictionary`](https://docs.godotengine.org/en/4.7/classes/class_dictionary.html) ("always passed by reference").

- `duplicate()` is shallow by default; `duplicate(true)` is deep (resources still shared).

## Collections

- Typed arrays `Array[int]`: element type checked on writes, `[]`, `[...] =`, and `+` (not on `push_back`);
  nested typed arrays are not supported. `Array` and `Array[Variant]` are the same type.
  [static_typing](https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/static_typing.html)
- Useful array methods: `map`, `filter`, `find` (−1 if none), `has` (or the `in` operator), `sort`
  (ascending, not stable), `duplicate`. Dictionary: `get(key, default)`, `has`, `keys`, `values`, `merge`.
  [`Array`](https://docs.godotengine.org/en/4.7/classes/class_array.html),
  [`Dictionary`](https://docs.godotengine.org/en/4.7/classes/class_dictionary.html)
- Packed arrays (e.g. `PackedInt32Array`) are faster/lower-memory but lack convenience methods like `map`.

## Error handling

- **No exceptions / try-catch** exists in GDScript (none is documented). Use return values and `is`/`as`
  guards, and report with `push_error()` / `push_warning()`
  ([`@GlobalScope`](https://docs.godotengine.org/en/4.7/classes/class_@globalscope.html)).
- `assert(condition, message)` runs **only in debug builds** and is ignored in release — keep it side-effect
  free. [basics](https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/gdscript_basics.html)
