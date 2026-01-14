package cmd

import (
	"fmt"
	"os"

	"ccm/internal/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ccm",
	Short: "Claude Code 供应商管理工具",
	Long: `ccm (Claude Code Manager) - 管理 Claude Code 的模型供应商

使用方法:
  ccm                   启动交互式 TUI 界面
  ccm list              列出所有供应商
  ccm add <name> --key  添加/配置供应商
  ccm run <name>        使用指定供应商启动 Claude Code

快捷键 (TUI):
  j/k, ↑/↓    移动选择
  Enter       使用选中的供应商启动 Claude
  e           编辑供应商配置
  d           设为默认
  t           测试连接
  /           搜索
  q           退出`,
	Run: func(cmd *cobra.Command, args []string) {
		// Launch TUI when no subcommand is provided
		result, err := ui.RunTUI()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// If user selected a provider to run, execute it
		if result != nil && result.RunProvider != "" {
			runCmd.Run(cmd, []string{result.RunProvider})
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
