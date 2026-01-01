package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ccm",
	Short: "Claude Code 供应商管理工具",
	Long: `ccm (Claude Code Manager) - 管理 Claude Code 的模型供应商

使用方法:
  ccm list              列出所有供应商
  ccm add <name> --key  添加/配置供应商
  ccm run <name>        使用指定供应商启动 Claude Code
  ccm generate          生成启动脚本

示例:
  ccm add doubao --key "sk-xxx"   配置豆包
  ccm run doubao                  启动 Claude Code（使用豆包）`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
