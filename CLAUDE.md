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

## Go Coding Standards

### Error Handling

- **Never ignore errors**: All `err` values must be checked
- **Wrap errors with context**: Use `fmt.Errorf("operation failed: %w", err)`
- **Return early on errors**: Avoid deep nesting

```go
// Good
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - never do this
_ = doSomething()
```

### Concurrency

- **Graceful shutdown**: All goroutines must exit via `context.Context` or channel signals
- **Avoid goroutine leaks**: Always ensure goroutines can terminate
- **Use sync primitives properly**: Protect shared state with mutexes or channels

```go
// Good
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case task := <-taskChan:
            process(task)
        }
    }
}
```

### Comments

- **Exported symbols**: All exported functions, types, and variables must have comments
- **Comment format**: Start with the symbol name (Go convention)

```go
// Provider represents an AI model provider configuration.
type Provider struct { ... }

// NewProvider creates a new Provider with the given options.
func NewProvider(opts ...Option) *Provider { ... }
```

### Project Structure

Follow standard Go project layout:

- `/cmd` - Main applications
- `/internal` - Private application code (not importable)
- `/pkg` - Public library code (if needed)
- `/docs` - Documentation (translations)

## Documentation & Language Policy

### Primary Language

- **All source code**: English (comments, variable names, commit messages)
- **All documentation in root**: English (README.md, CLAUDE.md, CONTRIBUTING.md)
- **Git commits**: English only

### Internationalization (i18n)

Translations are stored in `/docs` with language suffix:

```
docs/
├── README_zh-CN.md        # Chinese (Simplified)
├── CONTRIBUTING_zh-CN.md  # Chinese (Simplified)
├── README_ja.md           # Japanese (if needed)
└── ...
```

**Naming convention**: `<FILENAME>_<LANG-CODE>.md`

| Language | Code |
|----------|------|
| Chinese (Simplified) | `zh-CN` |
| Chinese (Traditional) | `zh-TW` |
| Japanese | `ja` |
| Korean | `ko` |

### Translation Guidelines

1. Keep translations in sync with English source
2. Do not translate code examples
3. Maintain consistent terminology across translations

## Safety Rules (--dangerously-allow-all-tools)

When `--dangerously-allow-all-tools` is enabled, strictly follow these guidelines:

### Git Pre-check

Before any multi-file refactoring or deletion:
- Verify Git working directory is clean (no uncommitted changes)
- Run `git status` to confirm

### Incremental Changes

- **Max 5 core files per change**: Never modify more than 5 business-critical files at once
- **Split complex refactors**: Break into smaller steps
- **Test after each step**: Run tests before proceeding

### Mandatory Testing

```bash
# After any code change
go test ./...
```

- **On test failure**: Immediately rollback or fix
- **Never continue**: Do not modify other files while tests are failing

### Forbidden Commands

- `rm -rf` and similar irreversible bulk deletions are **strictly prohibited**
- Use safe alternatives or manual confirmation for destructive operations

### Execution Logging

- Provide brief summary for each automated step
- Enable user to trace actions after the fact

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
