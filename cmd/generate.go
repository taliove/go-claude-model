package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"ccm/internal/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const scriptTemplate = `#!/usr/bin/env bash
# Claude Code - {{.DisplayName}}
# 由 ccm generate 自动生成

# 查找 claude 可执行文件
if [ -x "$HOME/claude-model/node_modules/.bin/claude" ]; then
    CLAUDE_BIN="$HOME/claude-model/node_modules/.bin/claude"
elif command -v claude &> /dev/null; then
    CLAUDE_BIN="$(command -v claude)"
else
    echo "错误: 未找到 claude 命令"
    echo "请先安装: npm install -g @anthropic-ai/claude-code"
    exit 1
fi

# 设置环境变量
export ANTHROPIC_AUTH_TOKEN="{{.APIKey}}"
export ANTHROPIC_BASE_URL="{{.BaseURL}}"
export ANTHROPIC_MODEL="{{.Model}}"
export API_TIMEOUT_MS=300000
export CLAUDE_CONFIG_DIR="$HOME/claude-model/configs/.claude-{{.Name}}"

# 确保配置目录存在
mkdir -p "$CLAUDE_CONFIG_DIR"

# 启动 claude
exec "$CLAUDE_BIN" "$@"
`

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "为已配置的供应商生成启动脚本",
	Long: `为所有已配置的供应商生成 shell 启动脚本

生成的脚本位于 ~/claude-model/bin/ 目录
将该目录加入 PATH 后，可直接使用 claude-<供应商名> 命令`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()

		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 加载配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		if len(cfg.Providers) == 0 {
			fmt.Println(red("没有已配置的供应商"))
			fmt.Println("请先使用 'ccm add <name> --key \"...\"' 配置供应商")
			return
		}

		// 创建 bin 目录
		home, _ := os.UserHomeDir()
		binDir := filepath.Join(home, "claude-model", "bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "%s 创建目录失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		// 解析模板
		tmpl, err := template.New("script").Parse(scriptTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 解析模板失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println(cyan("生成启动脚本:"))
		fmt.Println()

		count := 0
		for name, p := range cfg.Providers {
			if p.APIKey == "" {
				continue
			}

			scriptPath := filepath.Join(binDir, "claude-"+name)
			f, err := os.OpenFile(scriptPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				fmt.Printf("  %s %s: %v\n", red("✗"), name, err)
				continue
			}

			if err := tmpl.Execute(f, p); err != nil {
				f.Close()
				fmt.Printf("  %s %s: %v\n", red("✗"), name, err)
				continue
			}
			f.Close()

			fmt.Printf("  %s claude-%s\n", green("✓"), name)
			count++
		}

		fmt.Println()
		if count > 0 {
			fmt.Printf("已生成 %d 个脚本到 %s\n", count, binDir)
			fmt.Println()
			fmt.Println(cyan("将以下行添加到你的 ~/.bashrc 或 ~/.zshrc:"))
			fmt.Println()
			fmt.Printf("  %s\n", gray(`export PATH="$HOME/claude-model/bin:$PATH"`))
			fmt.Println()
			fmt.Println("然后重启终端或执行: source ~/.bashrc")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
