---
title: "cmdguard add"
status: planned
date: 2026-04-18
---

# cmdguard add

## Purpose

`cmdguard add` is intended to help append a new deny rule from CLI flags or
interactive prompts.

## Current status

This command is not yet fully specified for v1.

The main open questions are:

- whether the command should be purely flag-driven or interactive
- how much validation should happen before writing the file
- whether it should refuse to write a rule that lacks examples
- whether it should run `cmdguard test` automatically after mutation

## Minimum acceptable direction

If implemented in v1, `cmdguard add` should preserve the core safety principles
of the project:

- do not write invalid schema
- require or generate examples before completion
- make it obvious what changed in the config file

Until those details are fixed, this command should be treated as planned rather
than part of the stable v1 contract.
