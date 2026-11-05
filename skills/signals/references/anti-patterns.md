# Signal anti-patterns & alternatives

Companion to `SKILL.md`. These are **toolkit recommendations** informed by the documented API, not verbatim
doc rules; the API facts they rely on are cited in `SKILL.md` and the other references.

## Anti-patterns

- **Signals for synchronous request/response.** Signals are fire-and-forget notifications; they do not return
  a value to the emitter. If you need an answer, call a method (or pass a `Callable`) instead.
- **Connecting without a disconnect strategy.** Long-lived emitters connected to transient receivers leak
  dangling connections. Decide up front: `CONNECT_ONE_SHOT`, or disconnect in `_exit_tree()` / `_ExitTree()`.
  (In C#, `+=` and capturing lambdas specifically require manual `-=` — see `lifecycle-and-safety.md`.)
- **Reconnecting every frame / on re-entry without a guard.** Duplicate connects error
  (`ERR_INVALID_PARAMETER`). Guard with `is_connected()`.
- **Signals to issue commands downward.** Emitting a signal to tell a child what to do inverts ownership. Use
  a direct method call for imperative commands; reserve signals for "something happened" notifications.
- **Overusing a global event bus.** A single autoload that relays everything hides data flow and creates
  implicit coupling. Use it sparingly for genuinely cross-cutting events.
- **String-name typos.** Prefer the first-class `signal.connect(callable)` form (GDScript) and `SignalName`
  identifiers (C#) over raw strings.

## When another mechanism fits better

- **Parent → child command:** direct method call.
- **Need a return value:** method call, or pass a `Callable` the callee invokes.
- **Shared state many nodes read:** an autoload/singleton holding data (read by reference), not a signal.
- **Per-frame tight coupling:** direct calls in `_process`/`_physics_process` rather than signal spam.
