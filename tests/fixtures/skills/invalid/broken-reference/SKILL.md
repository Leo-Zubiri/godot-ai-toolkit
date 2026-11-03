---
name: broken-reference
description: Valid frontmatter, but the body links to a bundled file that does not exist.
---

# Invalid: broken reference

This links to [a missing reference](references/missing.md). The file does not exist, so the validator's
reference-integrity check must reject this skill.
