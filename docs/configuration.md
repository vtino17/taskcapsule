# Configuration

File: `.taskcapsule.json`

## Fields

| Field | Type | Description |
|-------|------|-------------|
| `version` | int | Schema version (currently 1) |
| `defaults.baseBranch` | string | Default base branch (default: `main`) |
| `defaults.branchPrefix` | string | Prefix for capsule branches (default: `task/`) |
| `defaults.gracefulShutdownSeconds` | int | Grace period before force kill (default: 5) |
| `defaults.healthTimeoutSeconds` | int | Health check timeout (default: 30) |
| `setup` | array | Commands to run during start |
| `services` | object | Named service configurations |
| `checks` | object | Named check configurations |

## Service configuration

| Field | Type | Description |
|-------|------|-------------|
| `command` | array | Command to run (required) |
| `workingDirectory` | string | Working directory relative to worktree |
| `environment` | object | Environment variables with `${PORT:name}` support |
| `inheritEnvironment` | array | Environment variable names to inherit |
| `health` | object | Health check configuration |

## Health check types

- `none` - No health check
- `process` - Process stays alive
- `tcp` - TCP port responds
- `http` - HTTP endpoint returns expected status

## Command security

Commands must be arrays, not shell strings. This prevents injection.
