extends Node

## Complete, parse-verifiable example: duplicate-connection safety, one-shot,
## and disconnecting on exit.

signal pinged

func _ready() -> void:
	_connect_once()
	# One-shot connections disconnect themselves after the first emission.
	pinged.connect(_on_pinged_once, CONNECT_ONE_SHOT)

func _connect_once() -> void:
	# Guard against the duplicate-connect error (ERR_INVALID_PARAMETER).
	if not pinged.is_connected(_on_pinged):
		pinged.connect(_on_pinged)

func _exit_tree() -> void:
	if pinged.is_connected(_on_pinged):
		pinged.disconnect(_on_pinged)

func _on_pinged() -> void:
	pass

func _on_pinged_once() -> void:
	pass
