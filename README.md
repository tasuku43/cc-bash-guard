# cmdguard

Declarative, testable command-string policy engine for AI agents and shells.

> **Status:** v1 core implementation in progress. See
> [`docs/dev/spec/README.md`](docs/dev/spec/README.md) for the v1.0 implementation contracts and
> [`docs/concepts/product-concept.md`](docs/concepts/product-concept.md) for the
> product concept.

## What it does

`cmdguard` is a tiny hook that decides whether a shell command is allowed to
run. It is called from Claude Code `PreToolUse`, `zsh` `preexec`,
`pre-commit`, CI, or anywhere else a command-string policy is useful.

Rules are declared in YAML. Every rule ships with block/allow examples, and
`cmdguard test` runs them as unit tests — so a rule change that would let
through a command it used to block fails CI, not production.

```yaml
# ~/.config/cmdguard/cmdguard.yml
version: 1
rules:
  - id: no-git-dash-c
    match:
      command: git
      args_contains:
        - "-C"
    message: "git -C is blocked. Change into the target directory and rerun the command."
    block_examples:
      - "git -C repos/foo status"
    allow_examples:
      - "git status"
      - "# git -C in comment"
```

Rules may use either:

- `match`: structured predicate matching, recommended for new rules
- `pattern`: a raw RE2 regular expression, kept as an escape hatch

## Non-goals

- LLM-assisted rule authoring and transcript mining live in a separate
  `cmdguard-claude-plugin` repository, so the core CLI has no LLM
  dependency.
- Non-`exec` action types (`write`, `fetch`, `mcp_call`) are post-v1.

See [`docs/README.md`](docs/README.md) for the current documentation map.

## Install

Not yet released. Once v1 ships:

```sh
brew install tasuku43/tap/cmdguard
# or
go install github.com/tasuku43/cmdguard/cmd/cmdguard@latest
```

## Setup

### 1. Initialize the user config

```sh
cmdguard init
```

This creates the default rule file at:

```text
~/.config/cmdguard/cmdguard.yml
```

### 2. Edit rules

Update `~/.config/cmdguard/cmdguard.yml` directly.
For every rule, keep both:

- `block_examples`
- `allow_examples`

### 3. Validate changes

Run the main authoring command after every rule edit:

```sh
cmdguard test
```

Use `cmdguard check` for spot checks against concrete commands:

```sh
cmdguard check --format json 'git -C repo status'
cmdguard check --format json 'AWS_PROFILE=read-only-profile aws s3 ls'
```

## Claude Code Hook Setup

Register `cmdguard eval` as a `PreToolUse` hook for `Bash`.

Example:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "cmdguard eval" }
        ]
      }
    ]
  }
}
```

### Hook ordering with other tools

If you also use another `PreToolUse` Bash hook such as `rtk hook claude`,
register `cmdguard eval` first.

Recommended order:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "cmdguard eval" }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "rtk hook claude" }
        ]
      }
    ]
  }
}
```

This keeps `cmdguard` as the first deny gate before other hook-side behavior
runs.

## License

MIT.
