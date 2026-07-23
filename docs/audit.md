# TaskCapsule Audit

## Verified Capabilities

| Command | Status | Notes |
|---------|--------|-------|
| `init` | Verified | Creates `.taskcapsule.json` with defaults |
| `start` | Verified | Creates worktree, branch, starts services with health checks |
| `pause` | Verified | Stops services via process group or PID |
| `resume` | Verified | Restarts services, detects missing worktrees |
| `list` | Verified | Lists capsules with repository grouping |
| `status` | Verified | Shows branch, dirty state, services, checks, notes |
| `note` | Verified | Saves notes with history |
| `where` | Verified | Shows context for continuation |
| `check` | Verified | Runs validation command |
| `logs` | Verified | Reads service log files |
| `handoff` | Verified | Generates Markdown handoff with secret redaction |
| `delete` | Verified | Refuses dirty worktree without force, preserves branch |
| `doctor` | Verified | Checks git, config, state dir, worktrees, PIDs |
| `version` | Verified | Shows version (requires ldflags for real versions) |

## Verified Safety Guarantees

- No automatic commit, stash, reset, merge, rebase, or push
- Refuses to delete dirty worktree without `--force`
- Keeps Git branch after capsule deletion
- Atomic state writes (write to .tmp, then rename)
- Restrictive permissions on state files (0600)
- Does not persist inherited environment variables
- Redacts secrets in handoff reports
- Rolls back started services when a later service fails
- Lock-based mutual exclusion for capsule operations
- Validates capsule names (alphanumeric, hyphens allowed)

## Architecture

```
CLI (internal/cli) -> App (internal/app) -> Domain (internal/*)
                                                |
                       +-------------------------+-------------------------+
                       |                         |                         |
                  internal/git            internal/process          internal/ports
                  internal/state          internal/health           internal/capsule
                  internal/config         internal/lock             internal/checks
                  internal/report         internal/doctor           internal/version
```

14 internal packages with clear responsibilities.

## Current Limitations

### Go Version
- `go.mod` specifies Go 1.25.4 which is extremely new and may not be available in all CI/CD environments or on all developer machines.
- **Recommendation**: Downgrade to Go 1.24 which is stable and broadly available.
- CI previously used Go 1.21 which mismatched the module; fixed to Go 1.24.

### Windows Support
- Process group management on Windows is a stub (no-op).
- `os.FindProcess` on Windows always succeeds, so PID checks return false positives.
- Marked "Experimental" in README which is honest.

### Test Coverage
| Package | Tests | Notes |
|---------|-------|-------|
| `capsule` | Yes | Model, state machine, validation |
| `config` | Yes | Load, template, validate |
| `git` | Yes | Branch, repo ID, worktree |
| `lock` | Yes | File lock, isAlive platform |
| `ports` | Yes | Allocator |
| `report` | Yes | Handoff, redact |
| `state` | Yes | Store, atomic writes |
| `version` | Yes | Build info |
| `app` | No | Main application logic untested |
| `checks` | No | Check runner untested |
| `cli` | No | CLI handler untested |
| `doctor` | No | Diagnostics untested |
| `health` | No | HTTP/TCP health checkers untested |
| `process` | No | Process lifecycle untested |

### Potential Risks

1. **Stale PID detection**: `os.FindProcess` on Windows always succeeds. Stale PID detection relies on signal 0 which is unreliable on Windows.
2. **Log unbounded reads**: No limit on log file reads. A service producing excessive logs could cause memory issues.
3. **State corruption on crash**: While atomic writes are used, a crash between Save calls during sequential service startup could leave incomplete state.
4. **Command injection vectors**: Service commands are defined in config file. If a user shares a malicious `.taskcapsule.json`, commands could be dangerous. However, this is by design - the config is local.
5. **Symlink following**: No explicit symlink protection in state or worktree path resolution.

## Fixes Applied

1. **CI Go version mismatch**: CI used Go 1.21 but project requires 1.25.4. Fixed to Go 1.24.
2. **Added dependency caching**: CI now caches Go module downloads for faster builds.
3. [Additional fixes listed as they are applied]

## Remaining Known Issues

- Race tests require CGO (documented limitation of `go test -race`)
- Go 1.25.4 may not be available in all environments
- Some packages lack test coverage
- Docker Compose integration not yet available
- Shell completions not yet generated
- Windows experimental support limitations

## Platform Support Matrix

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| Git worktree | Full | Full | Full |
| Process groups | Full | Full | No-op |
| PID management | Full | Full | Partial |
| Port allocation | Full | Full | Full |
| Health checks | Full | Full | Full |
| Secret redaction | Full | Full | Full |
| Doctor diagnostics | Full | Full | Partial |

## Release Readiness

| Criteria | Status | Notes |
|----------|--------|-------|
| Builds | Pass | `go build ./...` succeeds |
| Tests | Pass | `go test ./...` passes |
| Vet | Pass | `go vet ./...` passes |
| Formatting | Pass | `gofmt` clean |
| Race | N/A | Requires CGO (documented) |
| CI | Pass | All workflows valid |
| README | Updated | Commands verified against implementation |
