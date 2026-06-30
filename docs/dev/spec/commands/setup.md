---
title: "cc-bash-guard setup"
status: implemented
date: 2026-06-23
---

# cc-bash-guard setup

## Purpose

`cc-bash-guard setup` is a non-mutating first-run orientation command. It gives
users the shortest safe path from installation to a verified Claude Code hook
setup.

## Behavior

`setup` should:

- print the recommended first-run checklist
- show the user config path and Claude Code settings path
- point to profile comparison with `cc-bash-guard init --list-profiles --verbose`
- point to the policy authoring loop with `suggest`, `explain`, and `verify`
- point to `docs/user/START_HERE.md`

`setup` must not create files, edit Claude Code settings, write policy, or run
verification automatically.
