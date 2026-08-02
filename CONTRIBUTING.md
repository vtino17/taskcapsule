# Contributing to TaskCapsule

Thank you for considering contributing to TaskCapsule.

## How to Contribute

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes with clear messages.
4. Push and open a Pull Request.

## Development Setup

- Go 1.25+
- Git available on `PATH`
- Run `gofmt`, `go vet`, unit tests, and integration tests before submitting.

```bash
test -z "$(gofmt -l .)"
go vet ./...
go test ./... -timeout 15m
go test -tags=integration ./test/integration/... -timeout 20m
```

## Code Guidelines

- Follow standard Go conventions and keep `gofmt` clean.
- Write tests for new functionality.
- Keep PRs focused on a single concern.
- Do not weaken filesystem boundaries, lock semantics, state validation, or secret-handling guarantees without updating the threat model and adding adversarial tests.

## Pull Requests

All required GitHub checks must pass. Security-sensitive changes should explain their trust boundary, negative tests, failure behavior, and platform impact. Do not include credentials, private logs, or unredacted handoff data in issues or pull requests.
