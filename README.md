# TaskCapsule

TaskCapsule isolates coding tasks in Git worktrees, manages their local services, and preserves enough state to pause, resume, validate, and hand work off without losing context.

It is designed for developers and coding agents that work on several branches concurrently. TaskCapsule is a local orchestration tool; it is not a remote execution service, container sandbox, or security boundary for untrusted repositories.

## Status

The latest supported release is the newest version listed on the [GitHub Releases](https://github.com/vtino17/taskcapsule/releases) page. Linux and macOS are supported. Windows process management remains experimental and is tested separately from the Unix process-group implementation.

Release candidates must pass formatting, vet, unit, race, integration, cross-platform build, CodeQL, `govulncheck`, reproducible packaging, SBOM, and provenance gates. Passing those gates reduces known risk but does not guarantee the absence of defects.

## Install

TaskCapsule requires Git and Go 1.25 or newer. Use a currently supported Go patch release; older standard-library builds may contain known vulnerabilities.

```bash
go install github.com/vtino17/taskcapsule@latest
taskcapsule version
```

You can also download a checksum-listed archive from [GitHub Releases](https://github.com/vtino17/taskcapsule/releases). Published release archives include Linux amd64/arm64, macOS amd64/arm64, and Windows amd64 binaries.

## Quick start

Run these commands from a Git repository:

```bash
taskcapsule init
taskcapsule start my-feature --no-services
taskcapsule note my-feature "Parser implemented; add malformed-input tests next."
taskcapsule check my-feature -- go test ./...
taskcapsule pause my-feature
taskcapsule resume my-feature
taskcapsule handoff my-feature
```

Each capsule receives a managed worktree below `~/.taskcapsule/worktrees`. State, checks, and service logs are stored below `~/.taskcapsule`; Unix builds apply restrictive modes, while Windows storage inherits the current account's ACLs.

## Configuration

`taskcapsule init` creates `.taskcapsule.json`. Commands are arrays and are executed directly without shell parsing.

```json
{
  "version": 1,
  "defaults": {
    "baseBranch": "main",
    "branchPrefix": "task/",
    "gracefulShutdownSeconds": 5,
    "healthTimeoutSeconds": 30
  },
  "setup": [
    { "command": ["go", "mod", "download"] }
  ],
  "services": {
    "api": {
      "command": ["go", "run", "./cmd/api"],
      "workingDirectory": ".",
      "environment": {
        "APP_ENV": "development",
        "PORT": "${PORT:api}"
      },
      "inheritEnvironment": ["DATABASE_URL"],
      "health": {
        "type": "http",
        "url": "http://127.0.0.1:${PORT:api}/health",
        "expectedStatus": 200,
        "timeoutSeconds": 30
      }
    }
  },
  "checks": {
    "test": { "command": ["go", "test", "./..."] }
  }
}
```

Static, non-secret environment values belong in `environment`. Put only variable names in `inheritEnvironment`; their values are read from the TaskCapsule process environment and are not persisted to capsule state. Service working directories must remain inside the managed worktree, including after symlink resolution.

See [configuration](docs/configuration.md) for the complete schema.

## Commands

| Command | Purpose |
| --- | --- |
| `init` | Create `.taskcapsule.json` |
| `start` | Create a branch/worktree and start configured services |
| `pause` / `resume` | Stop and restart capsule services |
| `status` / `list` / `where` | Inspect capsule state and continuation context |
| `note` | Store a concise continuation note |
| `check` | Run an explicit validation command in the worktree |
| `logs` | Read bounded service-log tails |
| `handoff` | Generate a redacted Markdown handoff |
| `doctor` | Detect stale locks and inconsistent local state |
| `delete` | Remove a managed worktree after safety checks |
| `completion` | Generate shell completion |

Run `taskcapsule --help` for command syntax.

## Safety model

- Capsule, service, and check identifiers are validated before they can form filesystem paths.
- Destructive cleanup is restricted to the configured TaskCapsule worktree root.
- Loaded state must match the active repository and capsule before it is used.
- Service working-directory symlinks cannot escape the managed worktree.
- Child services receive a small baseline environment plus explicitly configured or inherited variables, rather than every parent secret.
- On Unix, state directories use mode `0700` and state, logs, check output, and handoff files use mode `0600`. Windows storage inherits the current account's ACLs; Windows support remains experimental.
- Operations use per-capsule exclusive locks and atomic state replacement.
- HTTP health checks use bounded timeouts. A configured HTTP health URL can make an outbound request, so configuration must be trusted.

TaskCapsule intentionally executes commands declared in the repository configuration. Review `.taskcapsule.json` before using TaskCapsule in an unfamiliar repository.

See [security architecture](docs/security.md), [architecture](docs/architecture.md), and the [security policy](SECURITY.md).

## Verification

```bash
go test ./... -timeout 15m
go test -race ./... -timeout 20m
go test -tags=integration ./test/integration/... -timeout 20m
go vet ./...
bash scripts/verify-reproducible-release.sh ./dist
```

The race detector may be unavailable on constrained kernels. GitHub CI runs the authoritative race gate on a supported Linux runner.

## Contributing and support

Read [CONTRIBUTING.md](CONTRIBUTING.md) before proposing a change. Use GitHub Issues for reproducible bugs and feature discussions. Report suspected vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## License

Apache License 2.0. See [LICENSE](LICENSE).
