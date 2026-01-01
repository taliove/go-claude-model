package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"ccm/internal/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <name>",
	Short: "使用指定供应商启动 Claude Code",
	Long: `使用指定供应商启动 Claude Code

示例:
  ccm run doubao       使用豆包启动
  ccm run deepseek     使用 DeepSeek 启动`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		red := color.New(color.FgRed).SprintFunc()

		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 加载配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		// 检查供应商是否已配置
		p, ok := cfg.Providers[name]
		if !ok || p.APIKey == "" {
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 未配置\n", red("错误:"), name)
			fmt.Fprintf(os.Stderr, "请先运行: ccm add %s --key \"你的API密钥\"\n", name)
			os.Exit(1)
		}

		// 查找 claude 可执行文件
		claudeBin := findClaudeBin()
		if claudeBin == "" {
			fmt.Fprintln(os.Stderr, red("错误: 未找到 claude 命令"))
			fmt.Fprintln(os.Stderr, "请先安装 Claude Code:")
			fmt.Fprintln(os.Stderr, "  npm install -g @anthropic-ai/claude-code")
			fmt.Fprintln(os.Stderr, "或在本地安装:")
			fmt.Fprintln(os.Stderr, "  cd ~/claude-model && npm install @anthropic-ai/claude-code")
			os.Exit(1)
		}

		// 设置环境变量
		os.Setenv("ANTHROPIC_AUTH_TOKEN", p.APIKey)
		os.Setenv("ANTHROPIC_BASE_URL", p.BaseURL)
		os.Setenv("ANTHROPIC_MODEL", p.Model)
		os.Setenv("API_TIMEOUT_MS", "300000")

		// 设置独立的配置目录
		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, "claude-model", "configs", ".claude-"+name)
		os.MkdirAll(configDir, 0755)
		os.Setenv("CLAUDE_CONFIG_DIR", configDir)

		// 获取剩余参数传递给 claude
		claudeArgs := cmd.Flags().Args()

		// 使用 syscall.Exec 替换当前进程
		err = syscall.Exec(claudeBin, append([]string{"claude"}, claudeArgs...), os.Environ())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 启动 claude 失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}
	},
}

// findClaudeBin 查找 claude 可执行文件
func findClaudeBin() string {
	home, _ := os.UserHomeDir()

	// 优先查找本地安装
	localBin := filepath.Join(home, "claude-model", "node_modules", ".bin", "claude")
	if _, err := os.Stat(localBin); err == nil {
		return localBin
	}

	// 查找全局安装
	path, err := exec.LookPath("claude")
	if err == nil {
		return path
	}

	return ""
}

func init() {
	rootCmd.AddCommand(runCmd)
}
