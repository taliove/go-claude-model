<div align="center">

# CCM

**Claude Code 多模型管理器**

轻松切换 Claude Code 的 AI 模型供应商

[![Go Version](https://img.shields.io/github/go-mod/go-version/taliove/go-claude-model)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/taliove/go-claude-model)](https://github.com/taliove/go-claude-model/releases)
[![License](https://img.shields.io/github/license/taliove/go-claude-model)](../LICENSE)

[English](../README.md) | **简体中文**

</div>

---

## 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/taliove/go-claude-model/main/scripts/install.sh | bash
```

## 快速开始

```bash
ccm init                        # 1. 启动引导
ccm add doubao --key "your-key" # 2. 添加供应商
ccm run doubao                  # 3. 启动 Claude Code
```

## 功能特性

| | 功能 | 说明 |
|---|------|------|
| ⚡ | **一键切换** | 快速切换不同供应商 |
| 🔐 | **安全存储** | API Key 安全存储，支持环境变量 |
| 🌐 | **多供应商** | 豆包、DeepSeek、通义千问、Kimi、GLM 等 |
| 📜 | **脚本生成** | 自动生成各供应商启动脚本 |
| 🔧 | **自定义供应商** | 支持任意 OpenAI 兼容 API |

## 支持的供应商

| 供应商 | 名称 | 默认模型 |
|--------|------|----------|
| `doubao` | 豆包（字节跳动） | doubao-seed-code-preview-latest |
| `deepseek` | DeepSeek | deepseek-chat |
| `qwen` | 通义千问（阿里） | qwen-plus |
| `kimi` | Kimi（月之暗面） | moonshot-v1-8k |
| `siliconflow` | 硅基流动 | deepseek-chat |
| `glm` | GLM（智谱AI） | glm-4 |
| `wanjie` | 万界 | - |

## 命令参考

| 命令 | 说明 |
|------|------|
| `ccm init` | 交互式设置向导 |
| `ccm list` | 列出所有已配置的供应商 |
| `ccm add <name> --key "key"` | 添加或配置供应商 |
| `ccm edit <name> --key "key"` | 更新供应商配置 |
| `ccm run <name>` | 使用指定供应商启动 Claude Code |
| `ccm switch` | 交互式切换供应商 |
| `ccm test <name>` | 测试供应商连接 |
| `ccm generate` | 生成启动脚本 |
| `ccm remove <name>` | 删除供应商 |

## 自定义供应商

```bash
ccm add custom --key "your-key" --url "https://api.example.com/v1" --model "gpt-4"
```

## 环境变量

支持通过环境变量设置 API Key（优先级高于配置文件）：

```bash
export CCM_API_KEY_DOUBAO="your-api-key"
ccm run doubao
```

## 其他安装方式

<details>
<summary>从源码安装</summary>

```bash
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model
make install              # 安装到 ~/.local/bin
# 或
sudo make install-global  # 安装到 /usr/local/bin
```

</details>

<details>
<summary>直接下载二进制</summary>

```bash
curl -L https://github.com/taliove/go-claude-model/releases/latest/download/ccm -o ccm
chmod +x ccm
sudo mv ccm /usr/local/bin/
```

</details>

## 卸载

```bash
make uninstall              # 本地安装
sudo make uninstall-global  # 全局安装
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

[MIT License](../LICENSE)
