# Semantic Schemas

Semantic matching lets policy match parsed command meaning instead of only raw
text. Semantic fields are command-specific, and the schema is selected by
`command.name`.

Inspect the installed registry:

```sh
cc-bash-guard help semantic
cc-bash-guard help semantic git
cc-bash-guard semantic-schema --format json
cc-bash-guard semantic-schema git --format json
```

`cc-bash-guard help semantic` lists commands from the semantic schema registry.
`cc-bash-guard help semantic <command>` shows fields, types, descriptions,
boolean notes, and examples for one command. `semantic-schema --format json`
prints the same registry as machine-readable JSON.

## Supported Commands

The current registry supports:

- `git`
- `aws`
- `kubectl`
- `gh`
- `helmfile`

Treat the CLI output as the source of truth for the installed binary. Commands
without a semantic schema should use `patterns`.

## Field Types

- `string`: exact value match.
- `[]string`: list values checked by the field-specific matcher.
- `bool`: `true` or `false` for parser-recognized command properties.

Unsupported semantic fields fail `verify` and include the supported fields for
that command. Unsupported value types fail `verify`; for example, `force:
"true"` is rejected because `force` must be a bool.

## Flags Fields

`flags_contains` and `flags_prefixes` match option tokens recognized by the
command-specific parser. They do not scan raw argv words, positional arguments,
or commands handled only by GenericParser fallback. Whether a token is available
for semantic flags therefore depends on the parser for the selected command.

## git

Use Git semantic fields for verbs, remotes, branches, refs, push force syntax,
reset mode, clean mode, diff staging, and parser-recognized flags.

Boolean fields:

- `force`: for `git push`, true only for `--force` or `-f`; for `git clean`,
  true for `-f` or `--force`.
- `force_with_lease`: for `git push`, true for `--force-with-lease`.
- `force_if_includes`: for `git push`, true for `--force-if-includes`.
- `hard`: true for `git reset --hard`.
- `recursive`: true for `git clean -d`.
- `include_ignored`: true for `git clean -x` or `--ignored`.
- `cached` and `staged`: true for `git diff --cached` or `--staged`.

To block every force-like push syntax, include all three force fields in deny
rules or write separate rules:

```yaml
permission:
  deny:
    - name: git destructive force push
      command:
        name: git
        semantic:
          verb: push
          force: true
    - name: git force-with-lease push
      command:
        name: git
        semantic:
          verb: push
          force_with_lease: true
    - name: git force-if-includes push
      command:
        name: git
        semantic:
          verb: push
          force_if_includes: true
```

## aws

Use AWS semantic fields for service, operation, profile, region, endpoint, dry
run, pager behavior, and parser-recognized flags.

Boolean fields:

- `dry_run`: true when `--dry-run` is present, false when `--no-dry-run` is
  present, and unset when neither form is recognized.
- `no_cli_pager`: true when `--no-cli-pager` is present.

```yaml
permission:
  allow:
    - name: AWS identity
      command:
        name: aws
        semantic:
          service: sts
          operation: get-caller-identity
      env:
        requires:
          - AWS_PROFILE
```

See also `docs/user/AWS_GUIDELINES.md`.

## kubectl

Use kubectl semantic fields for verb, resource, namespace, context, filename,
selector, container, dry run, force, recursion, and parser-recognized flags.

Boolean fields:

- `namespace_missing`: true when no namespace was selected.
- `all_namespaces`: true for `-A` or `--all-namespaces`.
- `selector_missing`: true when no selector was selected.
- `dry_run`: true for `--dry-run` or a dry-run value other than `none`; false
  for `--dry-run=none`; unset when absent.
- `force`: true for `--force`.
- `recursive`: true for `-R` or `--recursive`.

```yaml
permission:
  allow:
    - name: kubectl read-only
      command:
        name: kubectl
        semantic:
          verb_in:
            - get
            - describe
```

## gh

Use GitHub CLI semantic fields for `api`, `pr`, and `run` workflows, including
repository selectors, API method and endpoint, PR merge options, run rerun
options, and parser-recognized flags.

Boolean fields include:

- `web`: true for `-w` or `--web`.
- `paginate`, `input`, `silent`, `include_headers`: true for the corresponding
  `gh api` options.
- `draft`, `fill`: true for the corresponding `gh pr create` options.
- `force`: true for `gh pr checkout --force` or `-f`, and `gh run rerun
  --force`.
- `admin`, `auto`, `delete_branch`: true for the corresponding `gh pr merge`
  options.
- `failed`, `debug`, `exit_status`: true for the corresponding `gh run`
  options.

```yaml
permission:
  deny:
    - name: mutating GitHub API
      command:
        name: gh
        semantic:
          area: api
          method_in:
            - POST
            - PATCH
            - PUT
            - DELETE
```

## helmfile

Use helmfile semantic fields for verb, environment, file, namespace,
kube-context, selector, dry run, wait behavior, delete behavior, state values,
and parser-recognized flags.

Boolean fields include:

- `environment_missing`, `file_missing`, `namespace_missing`,
  `kube_context_missing`, `selector_missing`: true when the corresponding
  value was not selected.
- `interactive`: true for `--interactive`.
- `dry_run`: true for `--dry-run`.
- `wait`: true for `--wait`.
- `wait_for_jobs`: true for `--wait-for-jobs`.
- `skip_diff`: true for `--skip-diff`.
- `skip_needs`: true for `--skip-needs`.
- `include_needs`: true for `--include-needs`.
- `include_transitive_needs`: true for `--include-transitive-needs`.
- `purge`: true for `--purge`.
- `delete_wait`: true for `--delete-wait`.

```yaml
permission:
  ask:
    - name: production helmfile destroy
      command:
        name: helmfile
        semantic:
          verb: destroy
          environment: prod
```

## Unsupported Commands

When a command is not listed by `cc-bash-guard help semantic`, write a
`patterns` rule instead of `command.semantic`:

```yaml
permission:
  ask:
    - name: tool preview
      patterns:
        - "^my-tool\\s+preview(\\s|$)"
```
