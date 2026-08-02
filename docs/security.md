# Security architecture

## Trust boundary

TaskCapsule is a local orchestration tool. The operating-system account, Git executable, TaskCapsule binary, and `.taskcapsule.json` configuration are trusted. Commands in the configuration are intentionally executed with the current user's permissions; TaskCapsule is not a sandbox for untrusted repositories.

## Filesystem safety

- Capsule, service, and check names are restricted to bounded identifiers before they form paths.
- Repository IDs are fixed-length hashes.
- Loaded state must match the requested capsule and active repository.
- A loaded worktree path must remain below the configured TaskCapsule worktree root before it can be used or recursively removed.
- Service working directories are checked after symlink resolution and cannot escape the worktree.
- State is replaced through a same-directory temporary file, `fsync`, and rename.
- State directories use mode `0700`; state, service logs, check logs, and handoff files use mode `0600`.

## Command and environment handling

Commands are arrays passed directly to `exec.Command`; TaskCapsule does not interpolate them into a shell string. This removes shell parsing but does not make a configured command trustworthy.

Child services receive a minimal operating-system baseline environment. Static non-secret values are read from `environment`; explicitly named values are copied from the parent through `inheritEnvironment`. Unlisted parent variables are not forwarded. `${PORT:name}` placeholders are replaced after local port allocation.

Never store credentials directly in `.taskcapsule.json`. Use `inheritEnvironment` and an external secret provider.

## Process and concurrency safety

Lifecycle mutations use exclusive per-capsule lock files. Service processes are placed in process groups on Unix so rollback and pause operations can stop descendants. State transitions are persisted around service startup and shutdown, and partial startup triggers reverse-order rollback.

Windows process management remains experimental because its process-group semantics differ from Unix.

## Network behavior

TaskCapsule does not listen on a network socket. TCP and HTTP health checks connect only to endpoints explicitly configured by the user. Requests use bounded timeouts. Treat configuration from an unfamiliar repository as executable code and review health-check destinations before running it.

## Redaction

Generated handoff reports redact common bearer-token, authorization-header, password, API-key, secret, and private-key patterns. Redaction is defense in depth, not a guarantee. Review a handoff before sharing it outside its original trust boundary.

## Supply chain

CI actions are pinned to immutable commits. Pull requests run tests, the race detector, integration tests, cross-platform smoke tests, CodeQL, `govulncheck`, and OpenSSF Scorecard. Release builds create deterministic archives, SHA-256 checksums, a CycloneDX SBOM, and GitHub artifact attestations.
