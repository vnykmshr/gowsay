# Contributing to gowsay

## Getting Started

```bash
# Clone and build
git clone https://github.com/vnykmshr/gowsay.git
cd gowsay
make build

# Run tests
make test

# Run locally
./bin/gowsay "Hello World"
./bin/gowsay serve
```

## Making Changes

**Before submitting a PR:**
- Run `make audit` (runs tests, linting, security checks)
- Add tests for new functionality
- Update README if adding features or changing behavior

## Adding a New Cow

1. Edit `cow/cows.go`
2. Add cow name to `cowNames` slice
3. Add template to `cows` map in `init()`
4. Use `{{.Thoughts}}`, `{{.Eyes}}`, `{{.Tongue}}` placeholders
5. Test: `go run . -c yourcow "test message"`

## Code Style

- Use `gofmt` for formatting
- Run `go vet` and `staticcheck`
- Keep functions focused and testable
- Add comments for exported functions

## Pull Requests

- One feature/fix per PR
- Clear commit messages
- Reference issues if applicable
- PRs reviewed within a few days

## Questions?

Open an issue or discussion on GitHub.
