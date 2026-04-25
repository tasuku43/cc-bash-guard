---
title: "Security Test Matrix"
status: implemented
date: 2026-04-25
---

# Security Test Matrix

## 1. Scope

`cc-bash-guard` tests are security boundary tests, not only feature tests. The
regression matrix must keep the important allow/ask/deny boundaries stable
across refactors.

The matrix lives in table-driven Go tests and covers:

- compound command composition
- shell features that fail closed
- parser-specific and generic parser behavior
- raw and structured permission rules
- permission source merging with Claude settings
- rewrite and permission interaction

## 2. Required Assertions

Every security matrix case must assert:

- final outcome: `allow`, `ask`, or `deny`
- trace presence and the security-relevant trace step
- shell shape: `simple`, `compound`, or `unknown`
- shape flags when the boundary depends on shell structure

Trace assertions must include the security-relevant event, such as:

- `fail_closed`
- `composition`
- `composition.command`
- rewrite primitive names
- `claude_settings`
- `permission_sources_merge`

## 3. Invariants

The matrix must preserve these invariants:

- `deny` must not become `allow`
- unsafe shell shapes must not become automatic `allow`
- `patterns` allow must not bypass fail-closed shell safety
- parser removal or generic fallback must not widen a semantic rule to `allow`
- permission sources must merge as `deny > ask > allow > abstain`
- rewrite steps must not hide the final shell shape used for permission
  evaluation

## 4. Required Categories

The security matrix must include at least these categories.

### Compound

- `&&`
- `||`
- pipeline
- sequence
- nested compound expressions

Supported list and pipeline shapes may allow only when every extracted command
is individually allowed. Nested or mixed unsafe shapes must ask unless an
extracted command is denied.

### Shell Features

- subshell
- command substitution
- process substitution
- redirection

These features are unsafe for automatic allow. Deny rules still apply to
extracted commands when the parser can find them.

### Parser

- semantic parser, including the Git parser
- generic fallback
- unknown command

When a deny or ask rule requires semantic fields and the semantic parser is not
available, evaluation must ask instead of falling through to a broader allow.

### Permission

- `patterns` deny
- `patterns` ask
- `patterns` allow
- command/env predicates

`patterns` deny and ask keep their priority. `patterns` allow must remain
narrow and must not authorize unsafe shell syntax.

### Permission Source Merge

- `cc-bash-guard` abstain + Claude allow => allow
- `cc-bash-guard` ask + Claude allow => ask
- `cc-bash-guard` allow + Claude ask => ask
- `cc-bash-guard` allow + Claude deny => deny
- `cc-bash-guard` abstain + Claude deny => deny
- both sources abstain => final fallback ask

The merge matrix must assert that `deny` always wins, explicit `ask` is not
overridden by another source's `allow`, and no-match remains `abstain` until
the final source merge. E2E hook tests must assert trace distinguishes explicit
ask from fallback ask with the final merge reason
`all sources abstained; fallback ask`.

### Evaluation Normalization

- shell `-c` built-in evaluation
- absolute command path basename matching
- AWS profile semantic parsing

The matrix must assert the original command remains unchanged and the final
permission decision is correct. Trace entries must distinguish original tokens
from normalized command names where normalization is used for evaluation.
