# CLAUDE.md

CCM (Claude Code Manager) - A CLI tool for managing multiple AI model providers with Claude Code.

## Quick Commands

| Command | Description |
|---------|-------------|
| `make build` | Build for development |
| `make dev` | Build with debug symbols |
| `make verify` | Run all checks (fmt + lint + test) |
| `make install` | Install to ~/.local/bin |
| `make help` | Show all available targets |

## Development Workflow

### Daily Development

```bash
# 1. Build with debug symbols
make dev

# 2. Test locally
./bin/ccm list
./bin/ccm switch

# 3. Run all checks before commit
make verify
```

### Before Commit

```bash
make verify  # Runs: fmt + lint + test
```

### Release

```bash
./scripts/release.sh <version>
git push origin main v<version>
```

## Project Structure

```
ccm/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command
│   ├── add.go             # Add provider
│   ├── run.go             # Run Claude Code
│   ├── list.go            # List providers
│   ├── switch.go          # Interactive switch
│   └── ...
├── internal/
│   ├── config/            # Configuration management
│   ├── provider/          # Provider definitions
│   └── ui/                # UI components (promptui)
├── scripts/
│   ├── install.sh         # Installation script
│   └── release.sh         # Release script
├── Makefile               # Build automation
└── .goreleaser.yaml       # Release configuration
```

## Git Conventions

### Commit Message Format

```
<type>: <description>
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code refactoring |
| `test` | Test related |
| `docs` | Documentation |
| `chore` | Maintenance |

### Examples

```
feat: add new provider support
fix: resolve API key validation issue
docs: update README installation guide
```

## Architecture Notes

### Provider System

- Providers defined in `internal/provider/preset.go`
- Types: `TypeDirect` (official API) or `TypeProxy` (third-party)

### Configuration

- Config file: `~/.config/ccm/config.yaml`
- Environment variables: `CCM_<PROVIDER>_API_KEY`

### Interactive UI

- Uses `promptui` for arrow-key selection
- Masked input for API keys

## Quality Standards

### Before Every Commit

```bash
make verify
```

This runs:
1. `go fmt ./...` - Code formatting
2. `golangci-lint run` - Linting (must pass)
3. `go test -v -race ./...` - Tests with race detection

**Note: Lint failures will block releases.**

## Release Checklist

1. Run all checks: `make verify`
2. Validate config: `make check`
3. Release: `./scripts/release.sh <version>`
4. Push: `git push origin main v<version>`
