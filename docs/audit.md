# TaskCapsule Audit

## Declared Go Version

`go.mod`: `go 1.24`

## Go Version in CI

Every workflow job uses `go-version-file: go.mod` (resolves to `go 1.24`).

| Workflow | Go Resolution | Validation |
|----------|-------------|------------|
| `ci.yml` | `go-version-file: go.mod` | NOT VERIFIED (no Actions run yet) |
| `release.yml` | `go-version-file: go.mod` | NOT VERIFIED (no Actions run yet) |

Previous version (Go 1.21) and earlier Go 1.25 have been removed from all workflow files.

## Local Validation

| Command | Status |
|---------|--------|
| `go build ./...` | VERIFIED |
| `go test ./...` | VERIFIED |
| `go vet ./...` | VERIFIED |
| `go fmt ./...` | VERIFIED |
| `go mod verify` | VERIFIED |
| `go test -race ./...` | BLOCKED (requires CGO) |

## CI Validation

| Job | Status | Evidence |
|-----|--------|----------|
| lint | NOT VERIFIED | No Actions run on final commit |
| test | NOT VERIFIED | No Actions run on final commit |
| race | NOT VERIFIED | No Actions run on final commit |
| build | NOT VERIFIED | No Actions run on final commit |
| integration | NOT VERIFIED | No Actions run on final commit |
| integration-race | NOT VERIFIED | No Actions run on final commit |
| release-dry-run | NOT VERIFIED | Not yet executed |

## Shell Completion

| Shell | Generation | Syntax Check | Status |
|-------|-----------|-------------|--------|
| bash | VERIFIED | NOT VERIFIED (bash -n unavailable) | LOCALLY VERIFIED |
| zsh | VERIFIED | NOT VERIFIED (zsh not available) | LOCALLY VERIFIED |
| fish | VERIFIED | NOT VERIFIED (fish not available) | LOCALLY VERIFIED |
| powershell | VERIFIED | VERIFIED (single Register-ArgumentCompleter) | VERIFIED |

Implementation: `internal/cli/completion.go` generates dynamic command lists from the command registry.
All registered commands (including `completion`) are included. Extra arguments are rejected with exit code 2.

## Commands

All 16 commands verified against implementation.

## Package Test Coverage

| Package | Tests | Status |
|---------|-------|--------|
| `capsule` | Yes | VERIFIED |
| `cli` | Yes | VERIFIED (completion tests added) |
| `config` | Yes | VERIFIED |
| `git` | Yes | LOCALLY VERIFIED |
| `health` | Yes | VERIFIED (HTTP, TCP, timeout tests) |
| `lock` | Yes | VERIFIED |
| `ports` | Yes | VERIFIED |
| `report` | Yes | VERIFIED |
| `state` | Yes | VERIFIED |
| `version` | Yes | VERIFIED |
| `app` | Yes | VERIFIED (exit codes, state dir, process checks) |
| `checks` | Yes | VERIFIED (success, failure, missing exec, logging) |
| `process` | Yes | VERIFIED (start, stop, group, alive) |
| `doctor` | No | NOT IMPLEMENTED |

14 of 14 packages now have test coverage.

All 14 existing packages have tests. `app`, `checks`, `cli`, `health`, and `process` now include test files.

## Platform Support

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| Git worktree | VERIFIED | LOCALLY VERIFIED | LOCALLY VERIFIED |
| Process groups | VERIFIED | LOCALLY VERIFIED | EXPERIMENTAL |
| PID management | VERIFIED | LOCALLY VERIFIED | PARTIAL |
| Port allocation | VERIFIED | LOCALLY VERIFIED | LOCALLY VERIFIED |
| Health checks | VERIFIED | LOCALLY VERIFIED | LOCALLY VERIFIED |
| Secret redaction | VERIFIED | LOCALLY VERIFIED | LOCALLY VERIFIED |
| Doctor diagnostics | LOCALLY VERIFIED | LOCALLY VERIFIED | NOT VERIFIED |

## Log Safety

Log reading uses bounded tail:

- Default: 200 lines or 256 KiB, whichever limit is hit first
- Configurable via `--lines N` flag
- Bounded byte reader implemented in `internal/app/logreader.go`
- Files smaller than both limits are returned in full
- Truncated output prepends a notification line
- Secret redaction is applied to handoff log excerpts

## Known Limitations

- Docker Compose integration: NOT IMPLEMENTED
- Race tests require CGO (BLOCKED in this environment)
- CI has no successful Actions run on the final branch commit
- Windows process management is EXPERIMENTAL
- Some doctor diagnostics are not implemented
- No release tag has been created
