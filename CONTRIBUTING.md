# Contributing to cmdproxy

Thanks for contributing.

## Development Setup

Requirements:

- Go 1.25+
- Git

Common commands:

```sh
task build
task test
task vet
task vulncheck
task release:preflight
```

## Workflow

1. Read `AGENTS.md` and `docs/dev/backlog/README.md`.
2. Confirm or update the relevant spec in `docs/dev/spec/**` first.
3. Implement in small, reviewable slices.
4. Add tests, including non-happy paths where relevant.
5. Run quality checks before opening a PR.

## Required Quality Gate

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## Security-sensitive Changes

Treat the following changes as security-sensitive:

- hook handling and Claude Code integration
- rewrite primitives and rewrite flow control
- config loading, cache loading, and policy evaluation
- binary provenance and verification behavior

For those changes, include in the PR:

- the security impact in plain language
- the affected specs or backlog items
- test evidence
- any trust-model implications

## Pull Requests

Please include:

- problem statement and scope
- spec or backlog links
- test evidence
- any contract or behavior changes

## Release Model

- releases are tag-driven (`v*`) via GitHub Actions
- artifacts are published to GitHub Releases
- stable tags can update the Homebrew tap (`tasuku43/homebrew-cmdproxy`)
