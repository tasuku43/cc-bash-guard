# Troubleshooting

## Verified Artifact Missing Or Stale

Run:

```sh
cc-bash-guard verify
```

The hook fails closed when the verified artifact is missing or stale unless
`cc-bash-guard hook --auto-verify` is configured.

If your config uses `include`, every included YAML file is part of the verified
artifact. Editing an included policy or test file makes the artifact stale; run
`cc-bash-guard verify` again.

## Include Errors

Common include failures:

- the include entry is empty
- the path is a URL instead of a local file path
- the file is missing or is not a regular file
- a nested include forms a cycle

Relative include paths are resolved from the file that declares the include. For
example, if `policies/base.yml` includes `./git.yml`, that path resolves to
`policies/git.yml`.

## Unsupported Semantic Field

Inspect the registered fields for the command:

```sh
cc-bash-guard help semantic git
cc-bash-guard semantic-schema git --format json
```

If verify reports an unknown key, use the current permission shape:
`command`, `env`, and `patterns`.

## Command Without Semantic Support

Semantic matching only works for commands listed by:

```sh
cc-bash-guard help semantic
```

Use `patterns` for raw regex rules when a command has no semantic schema.

## Final Result Is ask

`abstain` means no matching rule. If all permission sources abstain,
`cc-bash-guard` falls back to `ask`.

Add an explicit `allow`, `ask`, or `deny` rule when you want a stable decision.

## Regex Pattern Not Matching

`patterns` match the original command string and parsed command elements,
including shell `-c` inner commands. Anchor patterns carefully:

```yaml
permission:
  allow:
    - name: pwd
      patterns:
        - "^pwd$"
```

In YAML double-quoted strings, escape backslashes. Single-quoted YAML strings
can be easier for complex regular expressions.

## AWS Profile Style

Prefer this style in project guidance:

```sh
AWS_PROFILE=myprof aws eks list-clusters
```

The AWS parser can still evaluate profile, service, and operation semantically.
See `docs/user/AWS_GUIDELINES.md`.

## Command Not Being Rewritten

`cc-bash-guard` evaluates commands but does not rewrite them. Parser-backed
normalization is evaluation-only.
