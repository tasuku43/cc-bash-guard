---
title: "cc-bash-guard semantic-schema"
status: implemented
date: 2026-06-23
---

# cc-bash-guard semantic-schema

## Purpose

`cc-bash-guard semantic-schema` exposes the installed binary's semantic match
registry for humans and tooling.

## Behavior

`semantic-schema` should:

- print all registered schemas with `--format json`
- print one command schema with `cc-bash-guard semantic-schema <command> --format json`
- reject unknown command names
- support `cc-bash-guard semantic-schema <command> --examples` for a compact
  human view with common fields, examples, and the explain/suggest/verify loop

The JSON output is the stable tooling surface. The `--examples` output is for
human orientation.
