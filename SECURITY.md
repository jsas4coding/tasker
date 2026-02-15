# Security

## Trust Model

Tasker reads `.tasker/config.yml` and `.tasker/tasks/*.yml` files and generates executable shell code in `Taskfile.yml` and `Makefile`. These generated files contain commands that run with the full privileges of the user executing them.

**Configuration files are trusted input.** Tasker treats all `.tasker/` files as developer-authored content. Never run `tasker generate` on configuration obtained from untrusted sources.

## What Tasker Does

- Reads YAML configuration from `.tasker/`
- Validates configuration against embedded JSON Schemas
- Validates identifiers (group keys, task keys, environment names) against strict patterns
- Escapes metadata fields (name, description) when interpolating into Makefile shell strings
- Generates `Taskfile.yml` and `Makefile` containing the commands defined in your configuration

## What Tasker Does NOT Do

- Execute any commands — it only generates files
- Sandbox or restrict the commands you define
- Validate that commands in `cmds:` are safe
- Load configuration from remote sources
- Make network requests

## Security Boundaries

### Validated by Schema + Runtime (defense in depth)

| Field | Pattern | Used In |
|-------|---------|---------|
| Group keys | `^[a-z][a-z0-9_-]*$` | Makefile targets, YAML keys |
| Task keys | `^[a-z][a-z0-9_-]*(:[a-z][a-z0-9_-]*)*$` | Makefile targets, YAML keys |
| Environment names | `^[a-z][a-z0-9_-]*$` | Shell preconditions |

### Escaped for Makefile Context

| Field | Treatment |
|-------|-----------|
| `name` | Shell-escaped (`"`, `$`, `` ` ``, `\`, newlines) |
| `description` | Shell-escaped in `@echo`, newline-sanitized in `##` comments |

### Passthrough by Design

| Field | Reason |
|-------|--------|
| `cmds` | Shell commands — the purpose of the tool |
| `dir` | Working directory for the task runner |
| `deps` | Task dependencies for the task runner |
| `vars` | Template variables for the task runner |
| `dotenv` | Dotenv file paths for the task runner |

These fields are written to the generated output as-is because they are intended to be consumed by the task runner (Taskfile) or Make. Sanitizing them would break functionality.

## Reporting

This project is provided as-is. There is no formal security reporting process. If you find an issue, you may fork the repository and apply fixes.
