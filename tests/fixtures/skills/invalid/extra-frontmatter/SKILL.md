---
name: extra-frontmatter
description: Valid name and description, but carries a disallowed frontmatter key.
license: MIT
---

# Invalid: disallowed frontmatter key

The frontmatter contract is `name` + `description` only. `license` (and any other
key, e.g. `compatibility`/`metadata`) must be rejected — that information belongs
in the skill manifest, not the frontmatter.
