# Tasker

Task bundler for [Taskfile.yml](https://taskfile.dev/) and Makefile generation.

Tasker reads structured configuration from `.tasker/config.yml` and `.tasker/tasks/*.yml`, then bundles everything into a single `Taskfile.yml` and `Makefile`. It keeps project roots clean, provides structured task navigation with `{group}:{environment}:{action}` naming, dotenv loading per environment, and environment guards that prevent running the wrong tasks in the wrong context.

## Why Tasker?

Task runners like [Task](https://taskfile.dev/) and Make are great, but as projects grow, a single `Taskfile.yml` or `Makefile` becomes hard to maintain. Tasker solves this by splitting task definitions into logical groups while generating a single output file that works with standard tooling.

- **Organized source** — tasks live in `.tasker/tasks/*.yml`, one file per group
- **Standard output** — generates `Taskfile.yml` and `Makefile` that work everywhere
- **Environment guards** — `ENV=prod` blocks dev tasks from running in production
- **Dotenv per environment** — different `.env` files for dev, test, and prod
- **Schema validation** — JSON Schemas catch config errors before generation
- **IDE support** — `yaml-language-server` directives for autocompletion in editors

## Installation

```bash
go install tasker.jsas.dev/cmd/tasker@latest
```

Or build from source:

```bash
git clone https://github.com/jsas4coding/tasker.git
cd tasker
go build -o bin/tasker ./cmd/tasker
```

## Quick Start

```bash
# Scaffold a new project (detects package managers automatically)
tasker init

# Edit your tasks
$EDITOR .tasker/config.yml
$EDITOR .tasker/tasks/*.yml

# Generate Taskfile.yml and Makefile
tasker generate

# Run tasks via Task or Make
task go:dev:build
make go-dev-build
```

## Project Structure

After `tasker init`, your project gets a `.tasker/` directory:

```
project/
├── .tasker/
│   ├── config.yml              # Main configuration (environments, groups, vars)
│   ├── tasks/                  # Task definitions (one file per group)
│   │   ├── go.yml
│   │   ├── lint.yml
│   │   └── test.yml
│   └── schemas/                # JSON Schemas (exported for IDE support)
│       ├── tasker.schema.json
│       └── tasks.schema.json
├── Taskfile.yml                # Generated — do not edit
└── Makefile                    # Generated — do not edit
```

The generated files (`Taskfile.yml` and `Makefile`) should be gitignored. They are regenerated from `.tasker/` source files with `tasker generate`.

## Configuration

### .tasker/config.yml

This is the main configuration file. It defines environments, variables, and task groups.

```yaml
# yaml-language-server: $schema=schemas/tasker.schema.json
name: my-project
description: My project tasks
version: "3"

environments:
  dev:
    name: Development
    description: Local development environment
    dotenv: [".env", ".env.local"]
  test:
    name: Testing
    description: Automated testing environment
    dotenv: [".env", ".env.test"]
  prod:
    name: Production
    description: Production environment
    dotenv: [".env", ".env.production"]

vars:
  PROJECT_NAME: my-project
  PROJECT_ROOT: "{{.ROOT_DIR}}"

groups:
  go:
    name: Go
    description: Go build and dependency management
  lint:
    name: Lint
    description: Code quality and linting
  test:
    name: Test
    description: Testing and coverage
```

**Fields:**

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Project display name (shown in help output) |
| `description` | No | Project description |
| `version` | No | Taskfile schema version (default: `"3"`) |
| `environments` | No | Environment definitions with dotenv configuration |
| `vars` | No | Global variables (supports [Taskfile template syntax](https://taskfile.dev/usage/#variables)) |
| `groups` | Yes | Task groups — each key maps to `.tasker/tasks/<key>.yml` |

### .tasker/tasks/*.yml

Each group has its own task file. The filename must match the group key in `config.yml` (lowercase).

```yaml
# yaml-language-server: $schema=../schemas/tasks.schema.json
tasks:
  dev:build:
    name: Build (dev)
    description: Build with race detector for development
    environment: dev
    cmds:
      - mkdir -p bin
      - go build -race -o bin/app ./cmd/app

  dev:run:
    name: Run (dev)
    description: Build and run in development mode
    environment: dev
    cmds:
      - go run ./cmd/app

  prod:build:
    name: Build (prod)
    description: Optimized production build
    environment: prod
    cmds:
      - mkdir -p bin
      - CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app ./cmd/app
```

**Task fields:**

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Display name for the task |
| `description` | No | Description shown in listings |
| `environment` | No | Restrict task to a specific environment |
| `cmds` | Yes | Shell commands to execute (in order) |
| `dir` | No | Working directory (supports Taskfile template syntax) |
| `deps` | No | Task dependencies (run before this task) |
| `silent` | No | Suppress command echoing (default: `false`) |

### Task Naming Convention

Task keys follow the format `{environment}:{action}`:

```yaml
tasks:
  dev:build:        # Development build
  dev:run:          # Development run
  dev:watch:        # Development file watcher
  prod:build:       # Production build
  prod:deploy:      # Production deployment
```

The full task key (used with `task` or `make`) includes the group prefix:

```
{group}:{environment}:{action}
  go:dev:build
  lint:dev:check
  test:dev:unit
```

For Make targets, colons become dashes: `go:dev:build` → `go-dev-build`.

## Environment Guards

Tasks that declare an `environment:` field get a precondition in the generated output. This prevents accidentally running dev tasks in production or vice versa.

```bash
# Works: running a dev task with ENV=dev (or without ENV)
task go:dev:build
ENV=dev task go:dev:build

# Blocked: ENV=prod prevents dev tasks from running
ENV=prod task go:dev:build
# → Task go:dev:build requires ENV=dev (current: prod)

# Works: prod tasks only run with ENV=prod
ENV=prod task go:prod:build
```

Tasks without an `environment:` field run regardless of the `ENV` value.

## Dotenv Loading

Each environment can declare dotenv files. The generated `Taskfile.yml` includes all declared dotenv files at the root level (deduplicated, ordered alphabetically by environment key).

```yaml
# In config.yml
environments:
  dev:
    dotenv: [".env", ".env.local"]
  prod:
    dotenv: [".env", ".env.production"]

# Generated Taskfile.yml will include:
# dotenv:
#   - .env
#   - .env.local
#   - .env.production
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `tasker` | Alias for `tasker generate` |
| `tasker generate` | Generate `Taskfile.yml` and `Makefile` from `.tasker/` |
| `tasker init` | Scaffold `.tasker/` directory with starter configuration |
| `tasker validate` | Validate configuration without generating output |
| `tasker list` | Display structured task list with groups and environments |
| `tasker completion <shell>` | Generate shell completions (`bash`, `zsh`, `fish`) |
| `tasker version` | Display version, build time, and git commit |

### tasker init

Detects package managers in the current directory and scaffolds a starter configuration:

| File Detected | Suggested Group |
|---------------|-----------------|
| `go.mod` | Go |
| `package.json` | Node.js (npm/pnpm/yarn) |
| `composer.json` | PHP (Composer) |
| `Cargo.toml` | Rust (Cargo) |
| `pyproject.toml` | Python |

If no package managers are detected, a `general` group is created.

### tasker list

Displays a structured view of all tasks, grouped by category:

```
Tasker - Task bundler for Taskfile.yml and Makefile generation

Environments:
  dev     Development          .env, .env.local
  test    Testing              .env, .env.test
  prod    Production           .env, .env.production

Go        Go build and dependency management
  go:dev:build          Build (dev)              Build with race detector
  go:dev:run            Run (dev)                Build and run in development
  go:prod:build         Build (prod)             Optimized production build

Lint      Code quality and linting
  lint:dev:check        Check All (dev)          Run fmt, lint, vet, and tests

Test      Testing and coverage
  test:dev:unit         Unit Tests (dev)         Run unit tests
  test:dev:coverage     Coverage (dev)            Run with coverage report
```

## Schema Validation

Tasker validates all configuration files against JSON Schemas before generating output. Schemas are embedded in the binary and also exported to `.tasker/schemas/` for IDE support.

### IDE Integration

YAML files include `yaml-language-server` directives for automatic validation and autocompletion in editors that support it (VS Code with YAML extension, IntelliJ, Neovim with LSP, etc.):

```yaml
# yaml-language-server: $schema=schemas/tasker.schema.json
```

### Schema URLs

Schemas are versioned and available on GitHub:

- `https://raw.githubusercontent.com/jsas4coding/tasker/v0.1.0/schemas/tasker.schema.json`
- `https://raw.githubusercontent.com/jsas4coding/tasker/v0.1.0/schemas/tasks.schema.json`

Replace `v0.1.0` with the desired version tag, or use `main` for the latest.

## Generated Output

### Taskfile.yml

The generated `Taskfile.yml` is a standard [Taskfile v3](https://taskfile.dev/) file. It includes:

- `version: "3"` — Taskfile schema version
- `dotenv:` — all declared dotenv files
- `vars:` — global variables from config
- A `default` task that runs `task --list`
- All tasks with `desc:`, `summary:`, `cmds:`, and optional `dir:`, `deps:`, `silent:`
- Environment guard `preconditions:` for tasks with an `environment:` field

### Makefile

The generated `Makefile` includes:

- `-include .env` and `export` for dotenv loading
- `.PHONY` declarations for all targets
- A `help` target (default) with structured output
- All tasks as targets with `## description` comments
- Environment guard shell checks
- Colons in task names replaced with dashes (`go:dev:build` → `go-dev-build`)

### Taskfile Version Compatibility

Tasker generates output compatible with [Taskfile v3](https://taskfile.dev/). The `version` field in `.tasker/config.yml` controls the Taskfile schema version in the generated output.

To use Taskfile, install it from [taskfile.dev/installation](https://taskfile.dev/installation/).

## Security

See [SECURITY.md](SECURITY.md) for the trust model and security boundaries.

Key points:
- Configuration files are trusted input — never run `tasker generate` on untrusted configs
- Identifier fields (group keys, task keys, environment names) are validated against strict patterns
- Metadata fields (name, description) are escaped when interpolated into Makefile shell strings
- Command fields (`cmds`, `dir`, `deps`, `vars`, `dotenv`) are passthrough by design

## License

[MIT](LICENSE) — provided as-is, no warranty.
