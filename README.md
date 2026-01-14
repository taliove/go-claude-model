<div align="center">

# CCM

**Claude Code Manager**

Seamlessly switch between AI model providers for Claude Code

[![Go Version](https://img.shields.io/github/go-mod/go-version/taliove/go-claude-model)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/taliove/go-claude-model)](https://github.com/taliove/go-claude-model/releases)
[![License](https://img.shields.io/github/license/taliove/go-claude-model)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/taliove/go-claude-model)](https://goreportcard.com/report/github.com/taliove/go-claude-model)
[![CI](https://github.com/taliove/go-claude-model/actions/workflows/ci.yml/badge.svg)](https://github.com/taliove/go-claude-model/actions)

**English** | [简体中文](docs/README_zh-CN.md)

</div>

---

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/taliove/go-claude-model/main/scripts/install.sh | bash
```

## Quick Start

```bash
ccm init                        # 1. Setup wizard
ccm add doubao --key "your-key" # 2. Add provider
ccm run doubao                  # 3. Launch Claude Code
```

## Features

| | Feature | Description |
|---|---------|-------------|
| ⚡ | **One-Click Switch** | Switch between providers instantly |
| 🔐 | **Secure Storage** | API keys stored safely with env var support |
| 🌐 | **Multi-Provider** | Doubao, DeepSeek, Qwen, Kimi, GLM, and more |
| 📜 | **Script Generation** | Auto-generate launch scripts for each provider |
| 🔧 | **Custom Providers** | Add any OpenAI-compatible API endpoint |

## Supported Providers

| Provider | Name | Default Model |
|----------|------|---------------|
| `doubao` | Doubao (ByteDance) | doubao-seed-code-preview-latest |
| `deepseek` | DeepSeek | deepseek-chat |
| `qwen` | Qwen (Alibaba) | qwen-plus |
| `kimi` | Kimi (Moonshot) | moonshot-v1-8k |
| `siliconflow` | SiliconFlow | deepseek-chat |
| `glm` | GLM (Zhipu AI) | glm-4 |
| `wanjie` | Wanjie | - |

## Commands

| Command | Description |
|---------|-------------|
| `ccm init` | Interactive setup wizard |
| `ccm list` | List all configured providers |
| `ccm add <name> --key "key"` | Add or configure a provider |
| `ccm edit <name> --key "key"` | Update provider configuration |
| `ccm run <name>` | Launch Claude Code with provider |
| `ccm switch` | Interactive provider switching |
| `ccm test <name>` | Test provider connection |
| `ccm generate` | Generate launch scripts |
| `ccm remove <name>` | Remove a provider |

## Custom Provider

```bash
ccm add custom --key "your-key" --url "https://api.example.com/v1" --model "gpt-4"
```

## Environment Variables

API keys can be set via environment variables (takes priority over config):

```bash
export CCM_API_KEY_DOUBAO="your-api-key"
ccm run doubao
```

## Alternative Installation

<details>
<summary>From Source</summary>

```bash
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model
make install              # Install to ~/.local/bin
# or
sudo make install-global  # Install to /usr/local/bin
```

</details>

<details>
<summary>Binary Download</summary>

```bash
curl -L https://github.com/taliove/go-claude-model/releases/latest/download/ccm -o ccm
chmod +x ccm
sudo mv ccm /usr/local/bin/
```

</details>

## Uninstall

```bash
make uninstall        # Local installation
sudo make uninstall-global  # Global installation
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

[MIT License](LICENSE)
