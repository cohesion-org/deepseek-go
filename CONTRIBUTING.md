# Contributing to deepseek-go

Thanks for your interest in contributing! This document explains how to get started.

## Code of Conduct

This project follows the [Go Community Code of Conduct](https://go.dev/conduct). Be respectful and constructive.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/deepseek-go.git`
3. Add the upstream remote: `git remote add upstream https://github.com/cohesion-org/deepseek-go.git`
4. Create a branch: `git checkout -b your-feature-name`

## Development Setup

Requires Go 1.26.0 or later.

```sh
go mod download
go build ./...
go test -v -short ./...     # offline tests (no API key needed)
```

## Code Style

- `goimports` or `gofmt` for formatting
- `golangci-lint run` for linting (config in `.golangci.yml`)
- Follow existing patterns for struct fields, JSON tags, and error handling
- Keep comments terse and focused on WHY, not WHAT
- Don't add unnecessary abstractions

## Testing

All changes must pass the existing test suite:

```sh
make test       # full test suite
make test-short # skip long-running tests
make test-race  # race detection
```

Tests use the built-in mock DeepSeek server (`internal/testutil/mock_api.go`) — no API key required for offline tests.

Live integration tests require `DEEPSEEK_LIVE_TESTS=1` and a valid API key:

```sh
DEEPSEEK_LIVE_TESTS=1 go test -v -tags=integration ./...
```

## Pull Request Process

1. Write a clear PR title (under 70 characters, imperative mood)
2. Reference any related issues in the description
3. Ensure `go build ./...`, `go vet ./...`, and `make test` pass
4. Add tests for new functionality
5. Update examples if adding new features
6. Keep PRs focused — one feature or fix per PR

## Commit Messages

Use the repo's convention:

```
area: brief description of change
```

Examples:

- `ci: bump Go version to 1.26`
- `deps: upgrade ollama to v0.24.0`
- `models: add DeepSeekV4Flash and DeepSeekV4Pro constants`

## Reporting Bugs

Use the [bug report template](https://github.com/cohesion-org/deepseek-go/issues/new?template=bug_report.md). Include:

- Go version (`go version`)
- deepseek-go version
- Steps to reproduce
- Expected vs actual behavior

## Feature Requests

Use the [feature request template](https://github.com/cohesion-org/deepseek-go/issues/new?template=feature_request.md). Describe the use case and why it matters.

## Getting Help

- [GitHub Discussions](https://github.com/cohesion-org/deepseek-go/discussions) — Q&A and community
- [Issue tracker](https://github.com/cohesion-org/deepseek-go/issues) — bugs and features

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
