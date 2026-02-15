# Tasker - CLAUDE.md

## Project Overview

Tasker is a Go CLI tool that reads structured configuration (`.tasker/config.yml` + `.tasker/tasks/*.yml`) and bundles into a single `Taskfile.yml` and `Makefile`.

## Architecture

```
tasker/
├── cmd/tasker/
│   └── main.go                      # Entry point
├── internal/
│   ├── cmd/                         # CLI commands (cobra)
│   │   ├── root.go                  # Root command, default → generate
│   │   ├── generate.go              # Bundle and output files
│   │   ├── init.go                  # Scaffold new project
│   │   ├── validate.go              # Config validation
│   │   ├── list.go                  # Structured help/list
│   │   ├── completion.go            # Shell completions for tasker CLI
│   │   └── version.go              # Version display (ldflags)
│   ├── config/                      # Configuration parsing
│   │   ├── config.go                # Config schema + loading
│   │   ├── group.go                 # Group + task schema, task file loading
│   │   ├── environment.go           # Environment schema
│   │   ├── detect.go                # Package manager detection
│   │   ├── schema.go                # JSON Schema validation (embedded)
│   │   └── schemas/                 # JSON Schema files (embedded copy)
│   ├── bundler/                     # Output generation
│   │   ├── taskfile.go              # Taskfile.yml generation
│   │   └── makefile.go              # Makefile generation
│   ├── resolver/
│   │   └── resolver.go              # Task resolution and env guards
│   ├── constants/
│   │   └── constants.go             # Shared constants (paths, permissions)
│   └── output/
│       └── output.go                # Styled terminal output
├── schemas/                         # JSON Schemas (canonical, versioned via GitHub)
│   ├── tasker.schema.json
│   └── tasks.schema.json
├── .tasker/                         # Self-referential config
│   ├── config.yml
│   ├── tasks/{go,lint,test}.yml
│   └── schemas/*.json
```

## Build & Run

```bash
go build -o bin/tasker ./cmd/tasker   # Build
tasker generate                       # Generate Taskfile.yml + Makefile
tasker validate                       # Validate config
tasker list                           # Show task list
tasker init                           # Scaffold new project
```

## Conventions

- Entry point: `cmd/tasker/main.go`, commands in `internal/cmd/`
- Task naming: `{environment}:{action}` inside task files
- Full task key: `{group}:{environment}:{action}`
- Task file naming: lowercase of group key (e.g., `go` → `go.yml`)
- All tasker config lives in `.tasker/` directory
- Environment guards via preconditions in generated Taskfile.yml
- Generated files (`Taskfile.yml`, `Makefile`) are gitignored
- JSON Schemas: canonical at `schemas/`, embedded copy at `internal/config/schemas/`, exported to `.tasker/schemas/` on init
- Schema `$id` uses versioned GitHub raw URLs
- When updating schemas, edit `schemas/*.json` first, then copy to `internal/config/schemas/`
- Version injection via ldflags at build time
- No magic numbers: use `internal/constants/`
- Styled output: use `internal/output/`

## Quality

- Linting: `golangci-lint` (golangci.yaml) + `revive` (revive.toml)
- Formatting: `goimports`
- All packages have godoc comments

## Language

- Code and docs: English
- User-facing messages: English
