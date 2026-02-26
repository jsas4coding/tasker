# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- `Tasker.json` output — JSON snapshot of fully resolved project (metadata, environments, groups, tasks)
- `tasker export` command — generates only `Tasker.json`
- Built-in `tasker:*` task group — auto-injected tasks (generate, validate, list, init, export, version) in all output artifacts
- Reserved group key validation — configs declaring `groups.tasker` are rejected
- `builtin` flag on groups and tasks in `Tasker.json` to distinguish injected vs user-defined

### Changed

- `tasker generate` now produces `Tasker.json` alongside `Taskfile.yml` and `Makefile`
- `tasker list` now shows built-in `tasker:*` commands in output

## [1.0.1] - 2026-02-15

### Added

- Initial implementation of Tasker CLI
- `.tasker/` directory structure: `config.yml`, `tasks/*.yml`, `schemas/`
- `tasker generate` command to bundle config into Taskfile.yml and Makefile
- `tasker init` command to scaffold project with package manager detection
- `tasker validate` command for config validation (structural + JSON Schema)
- `tasker list` command for structured task listing with groups and environments
- `tasker completion` command for shell completions (bash, zsh, fish)
- `tasker version` command with build-time version injection via ldflags
- Environment guards via preconditions (ENV variable check)
- Dotenv loading support per environment
- Package manager detection (Node.js, PHP, Go, Rust, Python)
- JSON Schema validation for config.yml and task files
- Schema export to `.tasker/schemas/` for IDE support
- Versioned JSON Schemas at `schemas/` with GitHub raw URLs
- Styled terminal output with color support (NO_COLOR aware)
- `internal/constants` package for shared constants
- `internal/output` package for styled terminal messages
- Linter configs: golangci-lint (golangci.yaml) + revive (revive.toml)
- EditorConfig and comprehensive .gitignore
- Self-referential `.tasker/` configuration for the project itself
- Claude Code skill (`/tasker`) for AI-assisted scaffolding
- MIT license
- SECURITY.md with trust model and security boundaries
- Makefile shell escaping for name/description fields
- Runtime identifier validation (defense in depth beyond JSON Schema)
- GitHub Actions: PR quality gate (lint, security, test, build, schema sync)
- GitHub Actions: release pipeline (cross-compile, checksums, schemas, GitHub Release)
