---
title: "Configuration Model"
status: implemented
date: 2026-04-18
---

# Configuration Model

## 1. Scope

This document defines where `cc-bash-guard` looks for configuration in v1.

## 2. Supported Locations

`cc-bash-guard` loads pure policy from two optional layers:

1. User-wide:
   - `$XDG_CONFIG_HOME/cc-bash-guard/cc-bash-guard.yml`, or
   - `~/.config/cc-bash-guard/cc-bash-guard.yml` by default
2. Project-local:
   - `<project-root>/.cc-bash-guard/cc-bash-guard.yml`
   - `<project-root>/.cc-bash-guard/cc-bash-guard.yaml`

The effective policy is the merge of:

- global `cc-bash-guard` policy
- project-local `cc-bash-guard` policy

Merge order is deterministic:

- user-wide policy is loaded first
- project-local policy is loaded second
- within each config file, included files are loaded before the including file
- rewrite rules append in source order
- permission rules append within each bucket in source order
- top-level E2E tests append in source order

`claude_permission_merge_mode` is no longer supported. If present, verification
fails. No configuration is required to choose permission merge behavior.

`cc-bash-guard` policy and Claude settings.json permissions are permission
sources. Each source returns `deny`, `ask`, `allow`, or `abstain`. `abstain`
means no matching rule. Sources are merged with one rule:

```text
deny > ask > allow > abstain
```

An explicit `ask` is not overridden by `allow` from another source. `deny`
always wins. The final fallback is `ask` only when all sources abstain.

Project root resolution is currently delegated to the Claude-aware runtime paths
used by `cc-bash-guard hook` and `cc-bash-guard verify`.

Missing files are allowed and treated as absent layers.

## 3. Includes

Config files may declare top-level `include`:

```yaml
include:
  - ./policies/git.yml
  - ./tests/git.yml
```

Include paths are local files only. Relative paths are resolved relative to the
file that declares the include. Included files may include other files.

Unsupported include shapes fail verification:

- empty entries
- URLs
- missing files
- paths that are not regular files
- include cycles
- shell expansion, environment expansion, command substitution, and globbing

The effective order is recursive include order, then the current file.
Permission bucket lists and top-level E2E tests concatenate in effective order.
No deep merge or ID-based override is currently supported.

## 4. Rule Identity

The current schema does not expose rule IDs. Rules are identified by their
position, source layer, bucket, selector, and effect in traces and validation
messages.

There is no ID-based override or collision behavior in the current contract.

## 5. Empty and Invalid States

- Missing file: allowed, treated as no configured rules
- Empty file: invalid configuration
- Invalid YAML: invalid configuration
- Valid YAML with schema errors: invalid configuration

Invalid configuration causes `cc-bash-guard hook` to return a deny response rather than
silently falling back to partial policy enforcement.

## 6. Future Extensions

These are still post-v1 concerns:

- rule packs
- explicit override semantics
- rule IDs and ID collision policy
- additional config layers such as repo-global or team-managed paths
