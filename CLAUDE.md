# CLAUDE.md

CCM (Claude Code Manager) - 管理 Claude Code 多模型供应商的 CLI 工具。

## 常用命令

```bash
go build -o ccm .     # 构建
go test ./...         # 测试
go fmt ./...          # 格式化
./ccm list            # 列出供应商
./ccm add <name> --key "key"  # 添加供应商
./ccm run <name>      # 启动 Claude
```

## 项目结构

- `cmd/` - CLI 命令实现
- `internal/config/` - 配置管理
- `internal/provider/` - 供应商定义

## Git 规范

提交信息格式：`<type>: <description>`

类型：
- `feat` - 新功能
- `fix` - 修复 bug
- `refactor` - 重构
- `test` - 测试相关
- `docs` - 文档更新

示例：`feat: 添加新供应商支持`

## 开发流程 (TDD)

1. 先写测试 (`*_test.go`)
2. 运行测试确认失败
3. 编写最小实现代码
4. 运行测试确认通过
5. 重构优化

## 发布流程

发布新版本前**必须**完成以下检查：

```bash
go fmt ./...              # 1. 代码格式化
golangci-lint run         # 2. Lint 检查（必须通过）
go test ./...             # 3. 所有测试必须通过
goreleaser check          # 4. 验证发布配置
```

发布命令：

```bash
./scripts/release.sh <version>  # 例如: ./scripts/release.sh 0.4.0
```

注意：
- 版本号遵循语义化版本 (SemVer)
- CI 流水线会自动验证 lint 和测试
- **lint 不通过会阻塞发布**
