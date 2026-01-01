package cmd

import (
	"fmt"
	"os"

	"ccm/internal/config"
	"ccm/internal/provider"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "显示供应商详细信息",
	Long:  "显示指定供应商的详细信息，包括配置状态、API URL、模型等",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		// 检查是否是预置供应商
		preset, isPreset := provider.Presets[name]
		if !isPreset {
			fmt.Fprintf(os.Stderr, "%s 未知供应商: %s\n", red("错误:"), name)
			fmt.Println(gray("使用 'ccm list' 查看所有供应商"))
			os.Exit(1)
		}

		// 获取用户配置
		cfg, _ := config.Load()
		p, hasProvider := cfg.Providers[name]
		configured := hasProvider && p.APIKey != ""

		status := green("已配置")
		if !configured {
			status = red("未配置")
		}

		fmt.Println()
		fmt.Printf(cyan("供应商详情: %s\n"), name)
		fmt.Println()

		fmt.Printf("  %s 显示名称:   %s\n", gray("├"), preset.DisplayName)
		fmt.Printf("  %s 状态:       %s\n", gray("├"), status)
		fmt.Printf("  %s 模型:       %s\n", gray("├"), yellow(getModelOrDefault(name, cfg)))
		fmt.Printf("  %s API URL:    %s\n", gray("├"), getBaseURLOrDefault(name, cfg))
		fmt.Printf("  %s 获取 Key:   %s\n", gray("└"), preset.KeyURL)

		if configured {
			fmt.Println()
			fmt.Printf("  %s 使用 %s 启动 Claude\n", gray("├"), green(fmt.Sprintf("ccm run %s", name)))
			fmt.Printf("  %s 使用 %s 更新模型\n", gray("└"), green(fmt.Sprintf("ccm add %s --model \"新模型\"", name)))
		} else {
			fmt.Println()
			fmt.Printf("  %s 使用 %s 配置此供应商\n", gray("└"), green(fmt.Sprintf("ccm add %s --key \"your-api-key\"", name)))
		}

		fmt.Println()
	},
}

func getModelOrDefault(name string, cfg *config.Config) string {
	if p, ok := cfg.Providers[name]; ok && p.Model != "" {
		return p.Model
	}
	return provider.Presets[name].Model
}

func getBaseURLOrDefault(name string, cfg *config.Config) string {
	if p, ok := cfg.Providers[name]; ok && p.BaseURL != "" {
		return p.BaseURL
	}
	return provider.Presets[name].BaseURL
}

func init() {
	rootCmd.AddCommand(showCmd)
}
