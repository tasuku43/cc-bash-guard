---
title: "Compatibility and Distribution"
status: implemented
date: 2026-04-18
---

# Compatibility and Distribution

## 1. Scope

This document defines the intended compatibility and distribution stance for
`cmdproxy` v2-style rule files.

## 2. Rule Schema Stability

The current rule schema is `version: 2`.

- New features must not silently change the meaning of valid v2 rule files
- Breaking schema changes require a new version number

## 3. Runtime Expectations

`cmdproxy` is intended to run as:

- a local CLI in developer environments
- a hook target for AI-agent and shell integrations
- a CI-time validation or enforcement command

The implementation should favor:

- static binary distribution
- predictable exit codes
- no runtime service dependency

## 4. Distribution Targets

Planned v1 distribution channels are:

- `go install github.com/tasuku43/cmdproxy/cmd/cmdproxy@latest`
- GitHub Releases
- Homebrew tap

Additional package managers are post-v1.

## 5. Platform Stance

v1 should target the major developer platforms used for CLI tooling:

- macOS
- Linux
- Windows

The exact release matrix is an implementation detail, but the user-facing docs
should describe the supported installation paths clearly.
