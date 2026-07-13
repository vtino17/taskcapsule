# Security

## Secret handling

TaskCapsule may store environment variable names in configuration, but NEVER their values.

When starting services:
1. Values are read from the TaskCapsule process environment.
2. Passed to child processes.
3. NOT written to capsule state.
4. NOT displayed in terminal output.
5. NOT included in handoff reports.

## Redaction

Handoff reports redact:
- Bearer tokens
- Authorization headers
- Passwords
- API keys
- Secrets
- Private keys

## Path security

Capsule names are validated as safe slugs: alphanumeric, dashes, underscores, max 64 chars. No path separators, no `.` or `..`.

## Command security

All commands are executed as arrays (not shell strings) via `exec.Command`. No shell injection possible.
