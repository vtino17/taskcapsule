# TaskCapsule Audit

## Declared Go Version

`go.mod`: `go 1.25.4`

## Tested Go Version

- Go 1.24: build, test, vet, fmt all PASS (locally verified)
- Go 1.25.4: build, test, vet, fmt all PASS (locally verified)

Minimum compatible version: Go 1.24 (verified with `GOTOOLCHAIN=local`).

## Go Version in CI

All jobs use `go-version-file: go.mod`, resolving to Go 1.25.4.

## CI Validation

| Job | Status | Evidence |
|-----|--------|----------|
| lint | NOT VERIFIED | No Actions run on final commit |
| test | NOT VERIFIED | No Actions run on final commit |
| race | NOT VERIFIED | Requires CGO |
| build | NOT VERIFIED | No Actions run on final commit |
| integration | NOT VERIFIED | Requires integration environment |
| integration-race | NOT VERIFIED | Requires integration environment |

## Local Validation

| Command | Status |
|---------|--------|
| `go build ./...` | VERIFIED |
| `go test ./...` | VERIFIED |
| `go vet ./...` | VERIFIED |
| `gofmt -l .` | VERIFIED |

Race tests: BLOCKED (requires CGO; not available in this environment)

## Shell Completion

| Shell | Generation | Syntax Check | Status |
|-------|-----------|-------------|--------|
| bash | VERIFIED | NOT VERIFIED (bash not available in this env) | LOCALLY VERIFIED |
| zsh | VERIFIED | NOT VERIFIED (zsh not available) | LOCALLY VERIFIED |
| fish | VERIFIED | NOT VERIFIED (fish not available) | LOCALLY VERIFIED |
| powershell | VERIFIED | VERIFIED (single Register-ArgumentCompleter) | VERIFIED |

Implementation: `internal/cli/completion.go` generates dynamic command lists from the command registry. Each shell gets correctly formatted output. Error handling for unknown shell returns exit code 2.

## Commands

| Command | Status | Notes |
|---------|--------|-------|
| `init` | VERIFIED | Creates `.taskcapsule.json` |
| `start` | VERIFIED | Worktree + branch + services with health checks |
| `pause` | VERIFIED | Stops via process group or PID |
| `resume` | VERIFIED | Restarts, detects missing worktrees |
| `list` | VERIFIED | With repository grouping |
| `status` | VERIFIED | Branch, dirty state, services, checks, note |
| `note` | VERIFIED | Saves notes with history |
| `where` | VERIFIED | Summary for continuation |
| `check` | VERIFIED | Runs validation command |
| `logs` | VERIFIED | Reads log files (truncation limit: NOT IMPLEMENTED) |
| `handoff` | VERIFIED | Markdown with secret redaction |
| `delete` | VERIFIED | Dirty worktree rejected without force; branch preserved |
| `doctor` | VERIFIED | Git, config, state, worktrees, PIDs |
| `completion` | VERIFIED | bash, zsh, fish, powershell |
| `version` | VERIFIED | Build info |

## Package Test Coverage

| Package | Tests | Status |
|---------|-------|--------|
| `capsule` | Yes | VERIFIED |
| `config` | Yes | VERIFIED |
| `git` | Yes | LOCALLY VERIFIED |
| `lock` | Yes | VERIFIED |
| `ports` | Yes | VERIFIED |
| `report` | Yes | VERIFIED |
| `state` | Yes | VERIFIED |
| `version` | Yes | VERIFIED |
| `app` | No | NOT IMPLEMENTED |
| `checks` | No | NOT IMPLEMENTED |
| `cli` | No | NOT IMPLEMENTED |
| `doctor` | No | NOT IMPLEMENTED |
| `health` | No | NOT IMPLEMENTED |
| `process` | No | NOT IMPLEMENTED |

## Platform Support

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| Git worktree | Full | Full | Full |
| Process groups | Full | Full | No-op (EXPERIMENTAL) |
| PID management | Full | Full | Partial |
| Port allocation | Full | Full | Full |
| Health checks | Full | Full | Full |
| Secret redaction | Full | Full | Full |
| Doctor diagnostics | Full | Full | Partial |

## Docker Compose

NOT IMPLEMENTED. No Docker Compose example or integration exists.

## Log Safety

NOT IMPLEMENTED. Log reads have no bounded limit. Large log files could cause memory issues.

## Security

- Secret redaction: VERIFIED (handoff reports redact API keys, tokens, passwords)
- State file permissions: 0600
- No network services listen by default
- No external API calls
- No automatic destructive Git operations
- Atomic state writes (write .tmp, rename)

## Remaining Limitations

- 6 packages lack test files
- Log reads are unbounded
- Race tests require CGO
- Docker Compose integration is not available
- Windows process management is a no-op stub
- Some doctor checks are missing (port occupancy, missing executables)
