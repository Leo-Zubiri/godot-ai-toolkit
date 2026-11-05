extends RefCounted

## Complete, parse-verifiable: defensive programming — GDScript has no try/catch,
## so guard with checks and report with push_error/push_warning; assert is debug-only.

func safe_divide(a: float, b: float) -> float:
	if is_zero_approx(b):
		push_error("safe_divide: division by zero")
		return 0.0
	return a / b

func require_object(value) -> bool:
	if value == null or not (value is Object):
		push_warning("require_object: expected an Object")
		return false
	return true

func checked_amount(amount: int) -> int:
	assert(amount >= 0, "amount must be non-negative")  # debug-only; stripped in release
	return amount
