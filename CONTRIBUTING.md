# Contributing to CCM

[中文版](docs/CONTRIBUTING_zh-CN.md)

Thank you for your interest in contributing to CCM!

## Development Setup

### Prerequisites

- Go 1.21+
- golangci-lint (optional, for local linting)
- goreleaser (optional, for local release testing)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model

# Install dev dependencies
make deps

# Build
make build

# Run
./bin/ccm --help
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Build binary |
| `make dev` | Build with debug symbols |
| `make test` | Run tests |
| `make lint` | Run linter |
| `make fmt` | Format code |
| `make verify` | Run all checks (fmt + lint + test) |
| `make install` | Install to ~/.local/bin |
| `make help` | Show all targets |

## Code Style

We follow standard Go conventions:

- Run `make fmt` before committing
- All code must pass `golangci-lint`
- Use meaningful variable and function names
- Keep functions focused and small

## Commit Guidelines

We use [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# New feature
git commit -m "feat: add new provider support"

# Bug fix
git commit -m "fix: resolve config loading issue"

# Refactoring
git commit -m "refactor: simplify config parser"

# Documentation
git commit -m "docs: update README"

# Tests
git commit -m "test: add config tests"
```

These messages are used to auto-generate the CHANGELOG.

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes
4. Run checks: `make verify`
5. Commit with conventional commit message
6. Push and create a Pull Request

## Release Process

### Before Release

```bash
# Run all checks
make verify

# Test release locally
make release-local
```

### Creating a Release

```bash
# Run release script
./scripts/release.sh <version>

# Push to trigger CI
git push origin main v<version>
```

GitHub Actions will automatically:
1. Run tests and lint
2. Build binaries for all platforms
3. Generate CHANGELOG
4. Create GitHub Release
5. Update Homebrew formula

### Local Release (Advanced)

```bash
# Set GitHub token
export GITHUB_TOKEN="your-token"

# Run GoReleaser
goreleaser release --clean
```

## Project Structure

```
.
├── cmd/                    # CLI commands (Cobra)
├── internal/
│   ├── config/             # Configuration management
│   ├── provider/           # Provider definitions
│   └── ui/                 # UI components (promptui)
├── scripts/
│   ├── release.sh          # Release script
│   └── install.sh          # Installation script
├── .github/workflows/
│   ├── ci.yml              # CI workflow
│   └── release.yml         # Release workflow
├── .goreleaser.yaml        # GoReleaser config
├── .golangci.yml           # Linter config
└── Makefile                # Build automation
```

## Adding New Providers

1. Edit `internal/provider/preset.go`
2. Add new provider to `Presets` map
3. Add provider name to `PresetOrder` slice
4. Run tests: `make test`
5. Commit: `git commit -m "feat(provider): add XXX provider"`
