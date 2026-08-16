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
| `environment` | object | Static, non-secret environment values with `${PORT:name}` support |
| `inheritEnvironment` | array | Parent environment variable names to copy without persisting their values |
| `health` | object | Health check configuration |

## Health check types

- `none` - No health check
- `process` - Process stays alive
- `tcp` - TCP port responds
- `http` - HTTP endpoint returns expected status

## Command security

Commands must be non-empty arrays and are passed directly to the operating system without shell-string interpolation. They still execute with the current user's permissions, so repository configuration must be trusted.

Service and check names are bounded identifiers containing letters, digits, `_`, or `-`. Environment variable names use the portable `[A-Za-z_][A-Za-z0-9_]*` form. Working directories must be relative and remain inside the capsule worktree after symlink resolution.

Do not store secrets in `environment`. List their variable names in `inheritEnvironment`; TaskCapsule reads those values at service start and does not write them to capsule state.
