package cmd

import (
	"fmt"
	"os"

	"ccm/internal/config"
	"ccm/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:     "switch",
	Aliases: []string{"sw"},
	Short:   "交互式切换供应商",
	Long: `交互式切换供应商

使用方向键选择供应商，输入关键字可搜索过滤
选择后直接启动 Claude Code`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 加载配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		// 构建供应商列表（包含未配置的）
		items := ui.BuildProviderItems(cfg, true)

		if len(items) == 0 {
			fmt.Println("尚未配置任何供应商")
			fmt.Println()
			fmt.Printf("运行 %s 开始配置\n", cyan("ccm init"))
			return
		}

		// 统计已配置数量
		configuredCount := 0
		for _, item := range items {
			if item.IsConfigured {
				configuredCount++
			}
		}

		fmt.Printf("已配置 %s 个供应商\n", green(fmt.Sprintf("%d/%d", configuredCount, len(items))))
		fmt.Println()

		// 使用箭头键选择
		selectedName, err := ui.SelectProvider(items, "选择要使用的供应商 (输入可搜索)")
		if err != nil {
			fmt.Fprintf(os.Stderr, "取消选择\n")
			return
		}

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
