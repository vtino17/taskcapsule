# Changelog

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
