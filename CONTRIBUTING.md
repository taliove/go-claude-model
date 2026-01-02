# Contributing to CCM

## 开发环境

### 依赖

- Go 1.21+
- golangci-lint (可选，用于本地 lint)
- goreleaser (可选，用于本地测试发布)

### 快速开始

```bash
# 克隆仓库
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model

# 编译
make build

# 运行
./ccm --help
```

## 常用命令

```bash
# 编译（输出到项目根目录）
make build

# 运行测试
make test

# 代码格式化
make fmt

# 运行 linter
make lint

# 整理依赖
make tidy

# 安装到 ~/.local/bin
make install

# 安装到 /usr/local/bin（需要 sudo）
sudo make install-global

# 清理构建产物
make clean

# 查看所有命令
make help
```

## 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 新功能
git commit -m "feat: add new provider support"
git commit -m "feat(config): add timeout option"

# Bug 修复
git commit -m "fix: resolve config loading issue"

# 性能优化
git commit -m "perf: optimize provider lookup"

# 重构
git commit -m "refactor: simplify config parser"

# 文档
git commit -m "docs: update README"

# 测试
git commit -m "test: add config tests"
```

这些提交信息会自动生成 CHANGELOG。

## 发布流程

### 准备发布

```bash
# 运行发布脚本（会自动执行预检查）
./scripts/release.sh 0.3.0

# 预检查包括：
# - 验证版本格式
# - 检查工作目录干净
# - 检查 tag 不存在
# - 运行测试
# - 运行 linter (如果安装了 golangci-lint)
# - 验证 GoReleaser 配置 (如果安装了 goreleaser)
```

### 自动化流程

推送 tag 后，GitHub Actions 会自动：

1. 再次运行测试和 lint（失败则阻止发布）
2. 使用 GoReleaser 构建 6 个平台的二进制
3. 生成 CHANGELOG
4. 创建 GitHub Release
5. 更新 Homebrew formula

### 本地测试发布

```bash
# 安装 goreleaser
go install github.com/goreleaser/goreleaser/v2@latest

# 本地测试构建（不会真正发布）
make release-local

# 验证 GoReleaser 配置
make check
```

### 手动发布（本地执行 GoReleaser）

如果你想跳过 GitHub Actions，在本地直接执行 GoReleaser 发布：

```bash
# 1. 确保代码已提交且工作目录干净
git status

# 2. 创建并推送 tag
git tag -a v0.3.0 -m "Release v0.3.0"
git push origin v0.3.0

# 3. 设置 GitHub Token（需要 repo 权限）
export GITHUB_TOKEN="your-github-token"

# 4. 执行发布
goreleaser release --clean

# 或者如果不想更新 Homebrew formula：
goreleaser release --clean --skip=homebrew
```

**注意事项：**
- 需要先创建 GitHub Personal Access Token (PAT)，权限需要包含 `repo`
- 如果要更新 Homebrew formula，需要对仓库有写权限
- 本地发布会跳过 GitHub Actions 的预检查，请确保已运行 `make test` 和 `make lint`

**快速发布脚本：**

```bash
# 一键发布（需要提前设置 GITHUB_TOKEN 环境变量）
./scripts/release.sh 0.3.0  # 创建 tag
GITHUB_TOKEN=xxx goreleaser release --clean  # 本地发布
```

## 项目结构

```
.
├── cmd/                    # CLI 命令实现
│   ├── root.go
│   ├── add.go
│   ├── list.go
│   ├── run.go
│   ├── generate.go
│   └── version.go
├── internal/
│   ├── config/             # 配置管理
│   ├── provider/           # Provider 定义
│   └── version/            # 版本信息
├── scripts/
│   ├── release.sh          # 发布脚本
│   └── install.sh          # 安装脚本
├── .github/workflows/
│   ├── ci.yml              # CI 工作流（PR/push）
│   └── release.yml         # 发布工作流（tag）
├── .goreleaser.yml         # GoReleaser 配置
├── .golangci.yml           # Linter 配置
└── Makefile                # 构建命令
```

## 添加新 Provider

1. 编辑 `internal/provider/provider.go`
2. 在 `Presets` map 中添加新的 provider 配置
3. 运行测试确保没有问题
4. 提交：`git commit -m "feat(provider): add XXX provider"`
