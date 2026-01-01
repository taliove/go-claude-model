# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CCM (Claude Code Manager) is a CLI tool written in Go for managing multiple AI model providers for Claude Code. It allows switching between providers (Doubao, DeepSeek, Qwen, Kimi, SiliconFlow, GLM) without manually modifying configuration files.

## Commands

```bash
# Build
go build -o ccm .

# Run tests
go test ./...

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy

# Run CLI
./ccm --help
./ccm list              # List all providers
./ccm add <name> --key "your-api-key"  # Configure a provider
./ccm run <name>        # Launch Claude with provider
./ccm generate          # Generate shell scripts
```

## Architecture

- **cmd/**: Cobra command implementations (add.go, generate.go, list.go, root.go, run.go)
- **internal/config/**: YAML configuration management (~/.claude-model/configs/providers.yaml)
- **internal/provider/**: Provider presets and struct definitions

### Key Patterns

1. **Provider Abstraction**: Providers have Name, DisplayName, APIKey, BaseURL, Model, KeyURL
2. **Environment Injection**: Sets `ANTHROPIC_AUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL` before exec
3. **Isolated Configs**: Each provider gets `~/.claude-model/configs/.claude-<provider>/` directory
4. **Preset Providers**: 6 pre-configured providers in internal/provider/provider.go

### Configuration

- User config: `~/.claude-model/configs/providers.yaml`
- Generated scripts: `~/claude-model/bin/claude-<provider>`
- Claude executable lookup: `~/claude-model/node_modules/.bin/claude` (local) or $PATH (global)
