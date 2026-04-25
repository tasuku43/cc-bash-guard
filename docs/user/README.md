# User Documentation

Start here if your goal is to use `cc-bash-proxy` in local workflows.

## Entry

- `docs/user/START_HERE.md`

## Current Focus

- user-wide config at `~/.config/cc-bash-proxy/cc-bash-proxy.yml`
- rule editing followed by `cc-bash-proxy verify`
- hook integration via `cc-bash-proxy hook`
- semantic match discovery with `cc-bash-proxy help semantic` and
  `cc-bash-proxy help semantic <command>`

## Semantic Rule Help

`match.semantic` is command-specific and selected by exact `match.command`.
Supported semantic commands are exposed by:

```sh
cc-bash-proxy help semantic
cc-bash-proxy help semantic git
cc-bash-proxy semantic-schema --format json
```

Use `args_contains` / `args_prefixes` for legacy raw-word matching, and
`semantic.flags_contains` / `semantic.flags_prefixes` for flags recognized by a
command-specific semantic parser.

## Intended Guide Set

- `RULES.md`: writing directive-based rules safely
- `CLAUDE_CODE.md`: Claude Code hook usage and permission layering
- `SHELL.md`: shell and CI integration patterns

These guides are not written yet, but this directory remains the intended home
for user-facing documentation.
