# Assurance record

This document describes reproducible evidence for TaskCapsule. It is not an external certification and does not guarantee defect-free operation.

## Required pull-request gates

| Control | Evidence |
| --- | --- |
| Formatting and static analysis | `gofmt`, `go vet` |
| Unit behavior | `go test ./...` plus coverage artifact |
| Data-race detection | `go test -race ./...` on a supported GitHub Linux runner |
| Lifecycle behavior | tagged integration suite, with a separate race run |
| Platform compatibility | macOS and Windows smoke jobs plus release cross-builds |
| Known Go vulnerabilities | pinned `govulncheck` module invocation |
| Static security analysis | CodeQL for Go |
| Workflow and repository posture | OpenSSF Scorecard SARIF |
| Release integrity | two-build reproducibility comparison and SHA-256 checksums |
| Release inventory | CycloneDX SBOM |
| Provenance | GitHub artifact attestation using OIDC |

## Local candidate validation

The hardening candidate is evaluated with the checksum-verified official Go 1.25.12 Linux arm64 toolchain. Go 1.24 was rejected after `govulncheck` identified reachable standard-library vulnerabilities fixed in supported Go 1.25 patches.

| Command | Result |
| --- | --- |
| `go test ./... -count=1 -timeout 15m` | Passed |
| `go test -tags=integration ./test/integration/... -timeout 20m` | Passed |
| `go vet ./...` | Passed |
| `go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...` | Passed: no vulnerabilities found |
| `bash scripts/verify-reproducible-release.sh <output>` | Passed: five target archives matched byte-for-byte and all checksums verified |
| `go test -race ./...` | Environment blocked: local kernel exposes an unsupported ThreadSanitizer VMA range |

The authoritative race result is therefore the protected GitHub CI job, which runs on a supported hosted runner.

## Security controls under test

- unsafe capsule, service, check, repository, and state path inputs fail closed;
- loaded state cannot redirect worktree cleanup outside the managed root;
- service working-directory symlinks cannot escape a worktree;
- unlisted parent environment values are not passed to child services;
- state replacement uses a unique same-directory temporary file, synchronization, and rename;
- on Unix, state directories are mode `0700` and state, service logs, check logs, and handoff files are mode `0600`; Windows storage inherits the current account's ACLs;
- action dependencies are pinned to immutable commits;
- published archives are compared byte-for-byte across two builds before release.

## Known limitations

- Windows process lifecycle management is experimental.
- Local ephemeral-port allocation has an unavoidable bind gap between reservation and service startup; services must handle a bind failure, and TaskCapsule rolls back partial startup.
- Commands and health destinations from `.taskcapsule.json` are trusted input and execute with the current user's authority.
- Redaction is pattern-based defense in depth and requires human inspection before a handoff crosses trust boundaries.
- A tag-triggered release with the new SBOM and attestation path is not considered verified until its GitHub run succeeds.
