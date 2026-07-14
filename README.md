# TaskCapsule

> **Pause one coding task. Handle the interruption. Resume without losing your place.**

[![CI](https://github.com/vtino17/taskcapsule/actions/workflows/ci.yml/badge.svg)](https://github.com/vtino17/taskcapsule/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/vtino17/taskcapsule?display_name=tag)](https://github.com/vtino17/taskcapsule/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/vtino17/taskcapsule)](https://go.dev/)
[![License](https://img.shields.io/github/license/vtino17/taskcapsule)](LICENSE)

TaskCapsule is a local-first CLI that turns a coding task into a resumable **capsule**: an isolated Git worktree plus its development processes, ports, logs, notes, and latest check result.

No cloud account. No daemon. No API key. No AI model. No automatic commit, stash, reset, or push.

```bash
taskcapsule start payment-timeout
taskcapsule note payment-timeout "Investigate duplicate retry"
taskcapsule pause payment-timeout

# Handle an urgent interruption
taskcapsule start urgent-hotfix

# Continue the original task later
taskcapsule resume payment-timeout
taskcapsule where payment-timeout
```

## The problem

Switching branches is easy. Reconstructing the whole task is not.

A normal interruption often means remembering:

- which worktree and branch belong to the task
- which frontend, API, or worker processes were running
- which ports they used
- which test failed last
- what you planned to do next

TaskCapsule manages that lifecycle as one named task:

```text
start → work → pause → switch → resume → handoff → delete safely
```

## Why not just use Git worktree, tmux, or Docker Compose?

Those tools remain useful. TaskCapsule coordinates the task-level lifecycle around them.

| Tool | Primary responsibility |
| --- | --- |
| Git worktree | Branch and working directory isolation |
| tmux | Terminal sessions |
| Docker Compose | Containerized services |
| direnv | Directory-specific environment variables |
| **TaskCapsule** | Worktree, local processes, ports, logs, notes, checks, and handoff as one task |

## Install

### Go

```bash
go install github.com/vtino17/taskcapsule/cmd/taskcapsule@latest
```

### Release binary

Download the binary for Linux, macOS, or Windows from [GitHub Releases](https://github.com/vtino17/taskcapsule/releases).

Linux and macOS are fully supported in v0.1. Windows process-tree management is experimental.

## Try it in 60 seconds

Inside an existing Git repository:

```bash
# Create the project configuration
taskcapsule init

# Review .taskcapsule.json, then create a capsule
taskcapsule start my-feature

# Save the thought you do not want to forget
taskcapsule note my-feature "Implement checkout validation next"

# Stop services while preserving the worktree and context
taskcapsule pause my-feature

# Restart the task later
taskcapsule resume my-feature
taskcapsule where my-feature
```

## Commands

| Command | Description |
| --- | --- |
| `init` | Create `.taskcapsule.json` |
| `start` | Create a capsule with an isolated worktree and services |
| `pause` | Stop services and release runtime resources |
| `resume` | Restart services and restore the task context |
| `list` | List capsules |
| `status` | Show detailed capsule state |
| `note` | Save the current thought or next action |
| `where` | Reconstruct where you left off |
| `check` | Run and record a validation command |
| `logs` | Read service logs |
| `handoff` | Generate a secret-safe Markdown handoff |
| `delete` | Remove a capsule and its worktree safely |
| `doctor` | Diagnose stale PIDs, missing worktrees, and local state |
| `version` | Show build information |

## Configuration

Create `.taskcapsule.json` with `taskcapsule init`.

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
      "environment": {
        "PORT": "${PORT:api}"
      },
      "health": {
        "type": "http",
        "url": "http://127.0.0.1:${PORT:api}/health"
      }
    }
  }
}
```

See [configuration documentation](docs/configuration.md) for the full schema.

## Safety by default

TaskCapsule is intentionally conservative around source code:

- never automatically commits, stashes, resets, merges, rebases, or pushes
- refuses to delete a dirty worktree unless `--force` is explicit
- keeps the Git branch after a capsule is deleted
- stores state atomically with restrictive permissions
- never persists inherited environment-variable values
- redacts likely secrets from generated handoff reports
- rolls back already-started services when a later service fails its health check

See [security documentation](docs/security.md).

## How it works

```text
CLI
 └─ Application lifecycle
     ├─ Git worktrees and branches
     ├─ Process groups and logs
     ├─ Dynamic ports and health checks
     ├─ Atomic capsule state and lifecycle locks
     └─ Notes, checks, diagnostics, and handoff reports
```

Each component lives under `internal/` behind focused interfaces. See [architecture documentation](docs/architecture.md).

## Current limitations

- Windows process-tree termination does not yet use Job Objects.
- Port allocation has a small listen-close-bind race because v0.1 has no daemon.
- Stale locks are detected but are not automatically repaired.
- HTTP health-check requests currently use a fixed per-request timeout.

## Feedback and contributions

Real workflow feedback is more valuable than feature guesses.

Open an issue and describe:

1. what interrupted your original task
2. which context was difficult to recover
3. which commands or tools you currently use

Contributions are welcome. Start with the architecture and testing docs, then check open issues for small, verifiable improvements.

## License

Apache 2.0
