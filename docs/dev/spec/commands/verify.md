---
title: "cc-bash-guard verify"
status: proposed
date: 2026-04-23
---

# cc-bash-guard verify

## Purpose

`cc-bash-guard verify` is a stricter trust-oriented check than `cc-bash-guard doctor`.

It exists to answer a narrower question:

**Can the current local `cc-bash-guard` setup be reasonably trusted as
part of the execution path?**

## Behavior

`cc-bash-guard verify` should:

- resolve and merge global and project-local `cc-bash-guard` policy
- resolve global and project-local settings for the target tool
- run the same config and rule validation used by `doctor`
- run rewrite tests, permission rule tests, and top-level E2E tests against the
  effective merged state
- compile and write a tool-specific verified hook artifact
- require build metadata to be visible in the current binary
- for Claude, fail if Claude Code settings exist but do not point at
  `cc-bash-guard hook`
- for Claude, fail if Claude Code settings use `cc-bash-guard hook` via PATH lookup
  rather than an absolute binary path
- for Claude, fail if an absolute Claude Code hook target does not exist or is not
  executable
- for Claude, fail if Claude Code points at a different `cc-bash-guard` binary than the one
  currently being verified

It should not require the target tool to be installed. If no tool settings file
is present, that condition should remain informational rather than fatal.

## Output

### Human-readable

The default output should include:

- the running version
- the target tool
- the visible VCS revision or an explicit missing marker
- the underlying doctor-style checks
- a final verified true/false result
- the artifact cache paths when verification also produced executable hook artifacts

### JSON

`cc-bash-guard verify --format json` should expose:

- `verified`
- `tool`
- `build_info`
- `report`
- `failures`
- `artifact_built`
- `artifact_cache`

## Relationship To `doctor`

- `doctor` is broad and diagnostic
- `verify` is narrow and trust-oriented

`doctor` may emit warnings that are acceptable in development. `verify` should
promote a smaller set of trust-critical conditions into failures.

## Hook Relationship

`cc-bash-guard hook` reads only verified artifacts at runtime.

- If a verified artifact exists and matches the current config hash, the hook uses it
- If the artifact is missing or stale, the hook returns a deny response with `invalid_config`
- The deny reason should tell the user to run `cc-bash-guard verify`
- `cc-bash-guard hook --auto-verify` opts in to the older implicit verify behavior:
  the hook tries verify once and retries loading the artifact
