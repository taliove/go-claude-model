package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ccm/internal/config"
	"ccm/internal/provider"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "交互式切换供应商",
	Long: `交互式切换供应商

显示所有供应商列表，选择后直接启动 Claude Code`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 加载配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		// 收集所有供应商
		allProviders := []string{}

		// 先添加预置供应商
		for name := range provider.Presets {
			allProviders = append(allProviders, name)
		}

		// 再添加用户自定义供应商
		for name := range cfg.Providers {
			if _, exists := provider.Presets[name]; !exists {
				allProviders = append(allProviders, name)
			}
		}

		if len(allProviders) == 0 {
			fmt.Println("尚未配置任何供应商")
			fmt.Println()
			fmt.Printf("运行 %s 开始配置\n", cyan("ccm init"))
			return
		}

		fmt.Println("请选择要使用的供应商:")
		fmt.Println()

		// 显示列表
		configuredCount := 0
		for i, name := range allProviders {
			var displayName string
			var isConfigured bool

			// 检查是否是预置供应商
			if p, exists := provider.Presets[name]; exists {
				displayName = p.DisplayName
			} else if p, exists := cfg.Providers[name]; exists {
				displayName = p.DisplayName
			}

			// 检查是否已配置
			if p, exists := cfg.Providers[name]; exists && p.APIKey != "" {
				isConfigured = true
				configuredCount++
			}

			// 检查环境变量
			if config.GetEnvAPIKey(name) != "" {
				isConfigured = true
				configuredCount++
			}

			if isConfigured {
				fmt.Printf("  %d. %s (%s) %s\n", i+1, cyan(name), displayName, green("✓"))
			} else {
				fmt.Printf("  %d. %s (%s) %s\n", i+1, cyan(name), displayName, yellow("✗ 未配置"))
			}
		}

		fmt.Println()
		fmt.Printf("已配置 %d/%d 个供应商\n", configuredCount, len(allProviders))
		fmt.Println()

		fmt.Print("请选择编号: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var selectedIdx int
		_, _ = fmt.Sscanf(input, "%d", &selectedIdx)

		if selectedIdx < 1 || selectedIdx > len(allProviders) {
			fmt.Println("无效的选择")
			os.Exit(1)
		}

		selectedName := allProviders[selectedIdx-1]

		// 检查是否已配置
		p, exists := cfg.Providers[selectedName]
		if !exists || (p.APIKey == "" && config.GetEnvAPIKey(selectedName) == "") {
			fmt.Printf("\n%s 供应商 '%s' 未配置\n", red("错误:"), selectedName)
			fmt.Printf("请先配置: ccm add %s --key \"your-api-key\"\n", selectedName)
			os.Exit(1)
		}

		fmt.Printf("\n正在启动 Claude Code (%s)...\n", selectedName)

		// 构建并执行 run 命令的参数
		runArgs := []string{"run", selectedName}

		// 直接调用 runCmd 的 Run 函数
		runCmd.Run(cmd, runArgs)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
