package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		if !ok {
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 未配置\n", red("错误:"), name)
			fmt.Fprintf(os.Stderr, "请先运行: ccm add %s --key \"你的API密钥\"\n", name)
			os.Exit(1)
		}

		// 获取 API Key (支持环境变量)
		apiKey := config.GetEffectiveAPIKey(name)
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 未设置 API Key\n", red("错误:"), name)
			fmt.Fprintf(os.Stderr, "请先运行: ccm add %s --key \"你的API密钥\"\n", name)
			fmt.Fprintf(os.Stderr, "或设置环境变量: export CCM_API_KEY_%s=\"your-key\"\n", strings.ToUpper(name))
			os.Exit(1)
		}

		// 检查 npm 是否安装
		if !hasNPM() {
			fmt.Fprintln(os.Stderr, red("错误: 未找到 npm 命令"))
			fmt.Fprintln(os.Stderr, "💡 解决方案:")
			fmt.Fprintln(os.Stderr, "   - macOS: brew install node")
			fmt.Fprintln(os.Stderr, "   - Ubuntu/Debian: sudo apt install npm")
			fmt.Fprintln(os.Stderr, "   - Fedora: sudo dnf install nodejs")
			os.Exit(1)
		}

		// 查找 claude 可执行文件
		claudeBin := findClaudeBin()
		if claudeBin == "" {
			fmt.Fprintln(os.Stderr, red("错误: 未找到 claude 命令"))
			fmt.Fprintln(os.Stderr, "💡 解决方案:")
			fmt.Fprintln(os.Stderr, "   全局安装: npm install -g @anthropic-ai/claude-code")
			fmt.Fprintln(os.Stderr, "   本地安装: cd ~/claude-model && npm install @anthropic-ai/claude-code")
			os.Exit(1)
		}

		// 设置环境变量
		os.Setenv("ANTHROPIC_AUTH_TOKEN", apiKey)
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

// hasNPM 检查 npm 是否已安装
func hasNPM() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

// getOSType 获取操作系统类型
func getOSType() string {
	return strings.ToLower(strings.SplitN(os.Getenv("OSTYPE"), ";", 2)[0])
}

func init() {
	rootCmd.AddCommand(runCmd)
}
