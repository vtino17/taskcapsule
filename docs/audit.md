# TaskCapsule Audit

## Declared Go Version

`go.mod`: `go 1.24`

## Go Version in CI

Every workflow job uses `go-version-file: go.mod` (resolves to `go 1.24`).

## CI Validation

Final PR commit: `f3b0f7d89edc6604c0ba8a77902b9cfde191c807`
Workflow run: [30006669553](https://github.com/vtino17/taskcapsule/actions/runs/30006669553)

| Job | Conclusion |
|-----|-----------|
| lint | CI VERIFIED / SUCCESS |
| test | CI VERIFIED / SUCCESS |
| race | CI VERIFIED / SUCCESS |
| build | CI VERIFIED / SUCCESS |
| integration | CI VERIFIED / SUCCESS |
| integration-race | CI VERIFIED / SUCCESS |
| release-dry-run | CI VERIFIED / SUCCESS |

## Local Validation

| Command | Status |
|---------|--------|
| `go build ./...` | VERIFIED |
| `go test ./...` | VERIFIED |
| `go vet ./...` | VERIFIED |
| `gofmt -l .` | VERIFIED |
| `go mod verify` | VERIFIED |
| Race tests (local) | BLOCKED (CGO unavailable) |

## Package Test Coverage

All 13 packages have tests:

| Package | Tests | Status |
|---------|-------|--------|
| `app` | Yes (exit codes, state dir, log reader) | VERIFIED |
| `capsule` | Yes (model, validation, state machine) | VERIFIED |
| `checks` | Yes (success, failure, missing exec) | VERIFIED |
| `cli` | Yes (completion, command dispatch) | VERIFIED |
| `config` | Yes (load, template, validate) | VERIFIED |
| `git` | Yes (branch, repo ID, worktree) | LOCALLY VERIFIED |
| `health` | Yes (HTTP, TCP, timeout, StatusError) | VERIFIED |
| `lock` | Yes (file lock, isAlive) | VERIFIED |
| `ports` | Yes (allocator) | VERIFIED |
| `process` | Yes (start, stop, group, helper pattern) | VERIFIED |
| `report` | Yes (handoff, redact) | VERIFIED |
| `state` | Yes (store, atomic writes) | VERIFIED |
| `version` | Yes (build info) | VERIFIED |

Note: The duplicate `internal/doctor` package was removed. It was dead code (no references to it existed anywhere in the codebase). All doctor functionality is provided by `internal/app.Doctor()`.

## Shell Completion

| Shell | Generation | Syntax Check | Status |
|-------|-----------|-------------|--------|
| bash | VERIFIED | NOT VERIFIED (bash -n unavailable) | LOCALLY VERIFIED |
| zsh | VERIFIED | NOT VERIFIED | LOCALLY VERIFIED |
| fish | VERIFIED | NOT VERIFIED | LOCALLY VERIFIED |
| powershell | VERIFIED | CI VERIFIED | VERIFIED |

## Log Safety

- Default: 200 lines / 256 KiB tail
- Configurable via `--lines N` flag
- Large files: tail from byte limit + truncation notice

## Platform Support

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| Git worktree | CI VERIFIED | CI VERIFIED | NOT VERIFIED |
| Process groups | CI VERIFIED | CI VERIFIED | EXPERIMENTAL |
| PID management | CI VERIFIED | CI VERIFIED | PARTIAL |
| Port allocation | CI VERIFIED | CI VERIFIED | CI VERIFIED |
| Health checks | CI VERIFIED | CI VERIFIED | CI VERIFIED |
| Secret redaction | CI VERIFIED | CI VERIFIED | CI VERIFIED |

## Tag-Triggered Release Publication

NOT VERIFIED. No release tag has been created.

## Known Limitations

- Windows process management is EXPERIMENTAL
- Docker Compose integration: NOT IMPLEMENTED
