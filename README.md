# TaskCapsule

> Pause one coding task and resume another without losing your place.

TaskCapsule is a CLI that lets developers save a coding task as a **capsule** — an isolated Git worktree with managed processes, logs, notes, and check results.

```bash
taskcapsule start payment-timeout
taskcapsule note payment-timeout "Investigate duplicate retry"
taskcapsule pause payment-timeout
taskcapsule resume payment-timeout
taskcapsule handoff payment-timeout
```

## Problem

Every task switch costs context: stopping servers, switching branches, finding the right file, remembering what to do next. TaskCapsule automates this so you can switch tasks in seconds.

## Installation

### From source

```bash
go install github.com/vtino17/taskcapsule/cmd/taskcapsule@latest
```

### From release

Download the latest binary from [GitHub Releases](https://github.com/vtino17/taskcapsule/releases).

## Quick start

```bash
# Initialize configuration
taskcapsule init

# Start a new task
taskcapsule start my-feature

# Save a note about what you're working on
taskcapsule note my-feature "Implementing checkout validation"

# Pause and switch to another task
taskcapsule pause my-feature

# Resume later
taskcapsule resume my-feature

# Generate a handoff report
taskcapsule handoff my-feature

# Delete when done
taskcapsule delete my-feature
```

## Commands

| Command     | Description                              |
| ----------- | ---------------------------------------- |
| `init`      | Create `.taskcapsule.json`              |
| `start`     | Create a new capsule with worktree + services |
| `pause`     | Stop all services and release resources  |
| `resume`    | Restart services from paused state       |
| `list`      | List all capsules                        |
| `status`    | Show detailed capsule state              |
| `note`      | Save a context note                      |
| `where`     | Show summary to continue working         |
| `check`     | Run a validation command in the worktree |
| `logs`      | Show service logs                        |
| `handoff`   | Generate a Markdown handoff report       |
| `delete`    | Remove a capsule and its worktree        |
| `doctor`    | Check TaskCapsule installation and state |
| `version`   | Show version information                 |

## Configuration

`.taskcapsule.json` example:

```json
{
  "version": 1,
  "defaults": {
    "baseBranch": "main",
    "branchPrefix": "task/"
  },
  "services": {
    "api": {
      "command": ["go", "run", "./cmd/api"],
      "health": {
        "type": "http",
        "url": "http://127.0.0.1:${PORT:api}/health"
      }
    }
  }
}
```

## Security

- TaskCapsule never stores environment variable values.
- Handoff reports redact secrets and include a security notice.
- Capsule state is stored with `0600` permissions.

## Platform support

| Platform | Status      |
| -------- | ----------- |
| Linux    | Full        |
| macOS    | Full        |
| Windows  | Experimental |

## Architecture

```
CLI -> Application Layer -> Git / Process / State / Health / Ports
```

Each component lives in `internal/` with clear interfaces.

## License

Apache 2.0
