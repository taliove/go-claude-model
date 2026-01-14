# 贡献指南

[English](../CONTRIBUTING.md)

感谢您对 CCM 项目的关注！

## 开发环境

### 依赖

- Go 1.21+
- golangci-lint（可选，用于本地 lint）
- goreleaser（可选，用于本地测试发布）

### 快速开始

```bash
# 克隆仓库
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model

# 安装开发依赖
make deps

# 编译
make build

# 运行
./bin/ccm --help
```

## 常用命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译 |
| `make dev` | 带调试符号编译 |
| `make test` | 运行测试 |
| `make lint` | 运行 linter |
| `make fmt` | 格式化代码 |
| `make verify` | 运行所有检查 (fmt + lint + test) |
| `make install` | 安装到 ~/.local/bin |
| `make help` | 显示所有命令 |

## 代码规范

遵循标准 Go 规范：

- 提交前运行 `make fmt`
- 所有代码必须通过 `golangci-lint`
- 使用有意义的变量和函数名
- 保持函数简洁

## 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 新功能
git commit -m "feat: add new provider support"

# Bug 修复
git commit -m "fix: resolve config loading issue"

# 重构
git commit -m "refactor: simplify config parser"

# 文档
git commit -m "docs: update README"

# 测试
git commit -m "test: add config tests"
```

这些提交信息会自动生成 CHANGELOG。

## Pull Request 流程

1. Fork 仓库
2. 创建功能分支：`git checkout -b feat/my-feature`
3. 进行修改
4. 运行检查：`make verify`
5. 使用规范的提交信息提交
6. 推送并创建 Pull Request

## 发布流程

### 发布前

```bash
# 运行所有检查
make verify

# 本地测试发布
make release-local
```

### 创建发布

```bash
# 运行发布脚本
./scripts/release.sh <version>

# 推送触发 CI
git push origin main v<version>
```

GitHub Actions 会自动：
1. 运行测试和 lint
2. 构建所有平台的二进制文件
3. 生成 CHANGELOG
4. 创建 GitHub Release
5. 更新 Homebrew formula

## 项目结构

```
.
├── cmd/                    # CLI 命令 (Cobra)
├── internal/
│   ├── config/             # 配置管理
│   ├── provider/           # Provider 定义
│   └── ui/                 # UI 组件 (promptui)
├── scripts/
│   ├── release.sh          # 发布脚本
│   └── install.sh          # 安装脚本
├── .github/workflows/
│   ├── ci.yml              # CI 工作流
│   └── release.yml         # 发布工作流
├── .goreleaser.yaml        # GoReleaser 配置
├── .golangci.yml           # Linter 配置
└── Makefile                # 构建自动化
```

## 添加新 Provider

1. 编辑 `internal/provider/preset.go`
2. 在 `Presets` map 中添加新 provider
3. 在 `PresetOrder` slice 中添加 provider 名称
4. 运行测试：`make test`
5. 提交：`git commit -m "feat(provider): add XXX provider"`
