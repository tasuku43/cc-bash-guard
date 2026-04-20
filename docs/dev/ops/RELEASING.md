# Releasing

## Overview

`cmdproxy` releases are tag-driven.

The intended pipeline is:

1. push a tag like `v0.1.0`
2. run CI-style preflight checks
3. build archives with GoReleaser
4. publish GitHub Release artifacts and `checksums.txt`
5. for stable tags without prerelease suffixes, optionally open a PR against
   `tasuku43/homebrew-cmdproxy`

## Release Inputs

- Workflow: `.github/workflows/release.yml`
- Packaging: `.goreleaser.yaml`
- Homebrew update script: `.github/scripts/update-homebrew-formula.sh`

## Artifacts

Every release should publish:

- platform archives
- `checksums.txt`

Checksums are part of the security story for `cmdproxy` because users are
trusting a binary that can rewrite commands before execution.

## Operator Checklist

1. Confirm CI is green on `main`.
2. Confirm `task release:preflight` passes locally if you are preparing the tag
   by hand.
3. Create and push a tag such as `v0.1.0`.
4. Confirm the GitHub Actions `Release` workflow succeeded.
5. Confirm the GitHub Release contains:
   - macOS archives for amd64 and arm64
   - Linux archives for amd64 and arm64
   - `checksums.txt`
6. Download one artifact and verify its checksum against `checksums.txt`.
7. Run the downloaded binary with:
   - `cmdproxy version --format json`
   - `cmdproxy verify --format json`
8. Confirm the reported VCS revision matches the intended release commit.
9. For stable tags, confirm a Homebrew formula PR was opened when Homebrew
   secrets are configured.

## Security Notes

- `checksums.txt` is the minimum release integrity signal and should always be
  present.
- Signed artifacts or attestations are still tracked separately in the security
  backlog.
- A release is not considered fully verified until an artifact has been
  downloaded and checked outside the CI environment.

## Homebrew Tap

Stable tags can update the Homebrew tap automatically when these secrets are
configured:

- `HOMEBREW_APP_ID`
- `HOMEBREW_APP_KEY`

Without those secrets, the GitHub Release still completes and Homebrew update
steps are skipped.
