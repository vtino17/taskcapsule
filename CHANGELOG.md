# Changelog

## [0.1.2] - 2026-07-24

### Fixed
- Scoped `doctor` branch checks to capsules belonging to the current repository
- Prevented false "branch no longer exists" warnings for capsules from other repositories
- Normalized repository path separators when computing fallback repository IDs on Windows
- Added isolated regression tests for cross-repository doctor diagnostics

## [0.1.1] - 2026-07-23

### Fixed
- Non-destructive process existence checks using signal 0
- Deterministic process tests without external `sleep` or `sh` dependencies
- Strict release archive and checksum validation
- Shell completion command coverage and deduplication
- Bounded service log reading

### Changed
- CI and release workflows use the Go version declared in `go.mod`
- Release packaging uses the shared `scripts/build-release-artifacts.sh`
- Removed the unused duplicate `internal/doctor` package

## [0.1.0] - 2026-07-14

### Added
- `init` - Create project configuration
- `start` - Create capsule with Git worktree and services
- `pause` - Stop all capsule services
- `resume` - Restart capsule services
- `list` - List all capsules in repository
- `status` - Show detailed capsule state
- `note` - Save context notes
- `where` - Show summary to continue working
- `check` - Run validation commands
- `logs` - View service logs
- `handoff` - Generate Markdown reports
- `delete` - Remove capsules
- `doctor` - Diagnose installation and state issues
- `version` - Show version
- Git worktree creation and removal
- Unix process group management
- Dynamic port allocation
- Health checks (none, process, TCP, HTTP)
- Atomic state store with locking
- Secret redaction in handoffs
- Cross-platform support (Linux, macOS, Windows experimental)
