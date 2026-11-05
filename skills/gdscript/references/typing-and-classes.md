# Typing, classes & version gates — details

Companion to `SKILL.md`. Cited to the Godot 4.7 docs; version labels verified against the per-version docs.

## Static typing

- `var x: int`, inference `var x := 5`, typed params `func f(a: int) -> void:`.
  [static_typing](https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/static_typing.html)
- Typed GDScript "improves performance by using optimized opcodes when operand/argument types are known at
  compile time" and "can detect more errors without even running the code" (4.7).
  **Version note:** the 4.0 page states this only as a *future* benefit — cite 4.7 for the present-tense claim.
- `as` casts (null on Object mismatch; errors for built-in-type mismatch); `is`/`is not` check runtime type.

## Classes

- `class_name X extends Y`; omitting `extends` ⇒ `RefCounted`. Inner classes: `class Name:` + `.new()`.
- `_init` is the constructor; a child must define `_init` if the base `_init` takes arguments; call the parent
  with `super(...)` or `super.method()`. No multiple inheritance.
  [basics](https://docs.godotengine.org/en/4.7/tutorials/scripting/gdscript/gdscript_basics.html)

## Version gates (do not use below the declared minimum)

| Feature | Introduced |
|---|---|
| Typed arrays `Array[T]` | 4.0 |
| `static var` / `@static_unload` | 4.1 |
| Typed `for x: T in …` | 4.2 |
| Typed dictionaries `Dictionary[K,V]` | 4.4 (per `gdscript_basics.html`) |
| `@abstract` classes/methods | 4.5 |

The `gdscript` skill's verified example files use only 4.0 features; newer syntax appears only as labeled
illustrative snippets in `SKILL.md`.
