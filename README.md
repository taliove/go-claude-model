# CCM (Claude Code Manager)

CCM 是一个用于管理 Claude Code 多模型供应商的命令行工具。通过 CCM，您可以轻松在不同 AI 模型供应商之间切换，无需手动修改配置。

## 功能特性

- **多供应商支持**: 支持豆包、DeepSeek、通义千问、Kimi、硅基流动、GLM 等主流 AI 模型提供商
- **快速切换**: 一键切换不同供应商，快速比较不同模型效果
- **配置管理**: 统一管理各供应商的 API Key 和配置
- **脚本生成**: 自动生成各供应商的启动脚本
- **安全存储**: API Key 安全存储，支持环境变量覆盖

## 安装

### 方式一: 从源码安装

```bash
# 克隆仓库
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model

# 本地安装 (推荐)
make install

# 或全局安装到 /usr/local/bin
sudo make install-global
```

### 方式二: 直接下载二进制

```bash
# 下载最新版本
curl -L https://github.com/taliove/go-claude-model/releases/latest/download/ccm -o ccm

# 添加执行权限
chmod +x ccm

# 移动到 PATH
sudo mv ccm /usr/local/bin/
```

## 使用方法

### 首次配置

```bash
# 启动引导
ccm init
```

### 基本命令

```bash
# 列出所有供应商
ccm list

# 添加/配置供应商
ccm add doubao --key "your-api-key"

# 更新供应商配置
ccm edit doubao --key "new-api-key"

# 测试供应商连接
ccm test doubao

# 启动 Claude Code
ccm run doubao

# 交互式切换供应商
ccm switch

# 删除供应商配置
ccm remove doubao

# 生成启动脚本
ccm generate
```

### 环境变量支持

支持通过环境变量设置 API Key (优先级高于配置文件):

```bash
export CCM_API_KEY_DOUBAO="your-api-key"
ccm run doubao
```

## 可用供应商

| 名称 | 显示名称 | 默认模型 |
|------|----------|----------|
| doubao | 豆包（字节跳动） | doubao-seed-code-preview-latest |
| deepseek | DeepSeek | deepseek-chat |
| qwen | 通义千问（阿里） | qwen-plus |
| kimi | Kimi（月之暗面） | moonshot-v1-8k |
| siliconflow | 硅基流动 | deepseek-chat |
| glm | GLM（智谱AI） | glm-4 |
| wanjie | 万界 | - |

## 自定义供应商

```bash
ccm add custom --key "xxx" --url "https://api.example.com/v1" --model "gpt-4"
```

## 生成启动脚本

```bash
ccm generate
```

生成的脚本位于 `~/claude-model/bin/`，可以这样使用:

```bash
# 直接运行脚本
~/claude-model/bin/claude-doubao

# 或添加到 PATH
export PATH="~/claude-model/bin:$PATH"
claude-doubao
```

## 卸载

```bash
# 本地安装卸载
make uninstall

# 全局安装卸载
sudo make uninstall-global
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
