package cmd

import (
	"fmt"
	"os"

	"ccm/internal/config"
	"ccm/internal/provider"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:   "default [name]",
	Short: "设置或显示默认供应商",
	Long: `设置或显示默认供应商

不带参数时显示当前默认供应商，带参数时设置默认供应商。
设置默认供应商后，运行 'ccm run' 时可以不指定供应商名称。

示例:
  ccm default           显示当前默认供应商
  ccm default doubao    设置 doubao 为默认供应商`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		// 无参数时显示当前默认
		if len(args) == 0 {
			defaultProvider := config.GetDefault()
			if defaultProvider == "" {
				fmt.Println("尚未设置默认供应商")
				fmt.Println()
				fmt.Printf("使用 %s 设置默认供应商\n", cyan("ccm default <name>"))
			} else {
				fmt.Printf("当前默认供应商: %s\n", green(defaultProvider))
			}
			return
		}

		name := args[0]

		// 验证供应商是否存在
		_, isPreset := provider.Presets[name]
		if !isPreset {
			cfg, _ := config.Load()
			if _, ok := cfg.Providers[name]; !ok {
				fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 不存在\n", red("错误:"), name)
				fmt.Fprintf(os.Stderr, "运行 %s 查看所有可用供应商\n", cyan("ccm list"))
				os.Exit(1)
			}
		}

		// 检查是否已配置
		if !config.IsConfigured(name) && config.GetEnvAPIKey(name) == "" {
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 尚未配置 API Key\n", yellow("警告:"), name)
			fmt.Fprintf(os.Stderr, "请先运行: %s\n", cyan(fmt.Sprintf("ccm add %s --key \"your-api-key\"", name)))
			fmt.Println()
		}

		// 设置默认供应商
		if err := config.SetDefault(name); err != nil {
			fmt.Fprintf(os.Stderr, "%s 设置默认供应商失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		fmt.Printf("%s 已设置 %s 为默认供应商\n", green("✓"), cyan(name))
		fmt.Println()
		fmt.Printf("现在可以直接运行 %s 启动 Claude Code\n", cyan("ccm run"))
	},
}

func init() {
	rootCmd.AddCommand(defaultCmd)
}
