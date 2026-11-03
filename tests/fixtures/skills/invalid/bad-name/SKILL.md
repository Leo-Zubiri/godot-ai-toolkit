---
name: Bad_Name
description: The name uses uppercase and an underscore, violating the [a-z0-9-] rule.
---

# Invalid: malformed name

`name:` is `Bad_Name`. Names must match `^[a-z0-9]+(-[a-z0-9]+)*$`. The validator must reject this
(it also will not equal the folder name `bad-name`).
