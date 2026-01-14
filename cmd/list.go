package cmd

import (
	"fmt"

	"ccm/internal/config"
	"ccm/internal/provider"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "列出所有供应商",
	Long:    "显示所有预置的供应商及其配置状态",
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()

		cfg, _ := config.Load()

		fmt.Println()
		fmt.Println(cyan("供应商列表:"))
		fmt.Println()

		// 显示默认供应商
		if cfg.Default != "" {
			fmt.Printf("  %s 默认: %s\n", yellow("★"), cyan(cfg.Default))
			fmt.Println()
		}

		for _, name := range provider.PresetOrder {
			preset := provider.Presets[name]
			p, hasProvider := cfg.Providers[name]
			configured := hasProvider && p.APIKey != ""

			status := red("✗")
			statusText := gray("未配置")
			if configured {
				status = green("✓")
				statusText = green("已配置")
			}

			// 供应商类型标签
			typeLabel := ""
			if preset.Type == provider.TypeProxy {
				typeLabel = magenta(" [代理]")
			}

			// 默认标记
			defaultMark := ""
			if cfg.Default == name {
				defaultMark = yellow(" ★")
			}

			fmt.Printf("  %s %-12s %s%s %s%s\n", status, name, preset.DisplayName, typeLabel, statusText, defaultMark)

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
		fmt.Printf("  %s\n", gray("ccm default <name>  # 设置默认供应商"))
		fmt.Printf("  %s\n", gray("ccm run             # 使用默认供应商启动"))
		fmt.Println()

		fmt.Println(yellow("管理命令:"))
		fmt.Printf("  %s\n", gray("ccm show <name>     # 查看供应商详情"))
		fmt.Printf("  %s\n", gray("ccm test <name>     # 测试 API 连接"))
		fmt.Printf("  %s\n", gray("ccm edit <name>     # 编辑供应商配置"))
		fmt.Printf("  %s\n", gray("ccm remove <name>   # 删除供应商"))
		fmt.Printf("  %s\n", gray("ccm switch          # 交互式切换供应商"))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
