# Changelog

## [Unreleased]

### Security

- Raise the minimum toolchain to Go 1.25 after vulnerability analysis found reachable issues in Go 1.24's standard library.
- Reject unsafe capsule, service, check, repository, and state path components.
- Validate loaded state and restrict recursive cleanup to the managed worktree root.
- Prevent working-directory symlinks from escaping a capsule worktree.
- Stop forwarding unlisted parent environment variables to child services.
- Store state, logs, checks, and handoffs with owner-only permissions.
- Replace state through unique synchronized temporary files.

### Added

- CodeQL, `govulncheck`, OpenSSF Scorecard, dependency updates, and cross-platform smoke gates.
- Deterministic release verification, CycloneDX SBOM generation, and artifact attestations.
- Regression coverage for path traversal, environment isolation, and symlink escape.

### Fixed

- Apply configured service environment values, inherited variables, dynamic ports, and working directories.
- Close parent service-log descriptors after process startup.
- Propagate check-log and state persistence failures.

## [0.1.2] - 2026-07-24

### Fixed

- Scoped `doctor` branch checks to capsules belonging to the current repository.
- Prevented false branch warnings for capsules from other repositories.
- Normalized repository path separators when computing fallback repository IDs on Windows.
- Added isolated regression tests for cross-repository doctor diagnostics.

## [0.1.1] - 2026-07-23

### Fixed

- Non-destructive process existence checks using signal 0.
- Deterministic process tests without external `sleep` or `sh` dependencies.
- Strict release archive and checksum validation.
- Shell completion command coverage and deduplication.
- Bounded service log reading.

### Changed

- CI and release workflows use the Go version declared in `go.mod`.
- Release packaging uses the shared `scripts/build-release-artifacts.sh`.
- Removed the unused duplicate `internal/doctor` package.

## [0.1.0] - 2026-07-14

### Added

- Capsule initialization, lifecycle, status, notes, checks, logs, handoff, deletion, diagnostics, and version commands.
- Git worktree creation and removal.
- Unix process-group management.
- Dynamic port allocation.
- Process, TCP, and HTTP health checks.
- Atomic state storage and per-capsule locking.
- Pattern-based handoff redaction.
- Linux and macOS support with experimental Windows support.
