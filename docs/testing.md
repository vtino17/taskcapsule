# Testing

## Unit tests

Cover config parsing, capsule validation, state transitions, port allocation, handoff generation, secret redaction, and more.

```bash
go test ./...
```

## Race detection

```bash
go test -race ./...
```

## Integration tests

Integration tests use temporary directories and Git repositories to verify the full lifecycle.

### Test scenarios

1. Happy path: start -> note -> check -> pause -> resume -> handoff -> delete
2. Dirty deletion is rejected without --force
3. Setup failure produces error state
4. Health check timeout cleans up services
5. Idempotent pause/resume commands
6. Secrets are never stored in state or handoff

## Running tests

```bash
go test ./... -count=1
go test -race ./...
go vet ./...
bash scripts/verify-reproducible-release.sh ./dist
```

GitHub CI additionally runs integration tests under the race detector, macOS and Windows smoke jobs, CodeQL, `govulncheck`, OpenSSF Scorecard, CycloneDX SBOM generation, and provenance attestation for published releases.
