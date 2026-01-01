package cmd

import (
	"fmt"

	"ccm/internal/config"
	"ccm/internal/provider"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有供应商",
	Long:  "显示所有预置的供应商及其配置状态",
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Println()
		fmt.Println(cyan("供应商列表:"))
		fmt.Println()

		for _, name := range provider.PresetOrder {
			preset := provider.Presets[name]
			cfg, _ := config.Load()
			p, hasProvider := cfg.Providers[name]
			configured := hasProvider && p.APIKey != ""

			status := red("✗")
			statusText := gray("未配置")
			if configured {
				status = green("✓")
				statusText = green("已配置")
			}

			fmt.Printf("  %s %-12s %s %s\n", status, name, preset.DisplayName, statusText)

			if configured {
				// 显示已配置详情
				fmt.Printf("    %s 模型: %s\n", gray("├"), yellow(p.Model))
				fmt.Printf("    %s URL: %s\n", gray("├"), p.BaseURL)
			}
			fmt.Printf("    %s 获取 API Key: %s\n", gray("└"), gray(preset.KeyURL))
			fmt.Println()
		}

		fmt.Println(yellow("快速开始:"))
		fmt.Printf("  %s\n", gray("ccm add <name> --key \"your-api-key\""))
		fmt.Printf("  %s\n", gray("ccm run <name>"))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
