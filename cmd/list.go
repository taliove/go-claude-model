package cmd

import (
	"fmt"
	"os"

	"ccm/internal/config"
	"ccm/internal/provider"
	"ccm/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var interactiveMode bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "列出所有供应商",
	Long: `显示所有预置的供应商及其配置状态

使用 -i 标志进入交互模式，可直接选择操作`,
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

		// 交互模式
		if interactiveMode {
			runInteractiveMode(cfg)
			return
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

func runInteractiveMode(cfg *config.Config) {
	red := color.New(color.FgRed).SprintFunc()

	actions := []string{
		"启动供应商 (run)",
		"配置供应商 (add)",
		"设为默认 (default)",
		"测试连接 (test)",
		"退出",
	}

	idx, err := ui.SelectAction(actions, "选择操作")
	if err != nil || idx == 4 {
		return
	}

	// 根据操作选择供应商
	var includeUnconfigured bool
	var label string

	switch idx {
	case 0: // run - 只显示已配置的
		includeUnconfigured = false
		label = "选择要启动的供应商"
	case 1: // add - 显示所有
		includeUnconfigured = true
		label = "选择要配置的供应商"
	case 2: // default - 只显示已配置的
		includeUnconfigured = false
		label = "选择要设为默认的供应商"
	case 3: // test - 只显示已配置的
		includeUnconfigured = false
		label = "选择要测试的供应商"
	}

	items := ui.BuildProviderItems(cfg, includeUnconfigured)

	if len(items) == 0 {
		fmt.Println("没有已配置的供应商，请先运行 ccm add <name> --key \"xxx\"")
		return
	}

	selectedName, err := ui.SelectProvider(items, label)
	if err != nil {
		return
	}

	// 执行对应操作
	switch idx {
	case 0: // run
		fmt.Printf("\n正在启动 Claude Code (%s)...\n", selectedName)
		runCmd.Run(nil, []string{selectedName})
	case 1: // add
		fmt.Println()
		p := provider.Presets[selectedName]
		if p.KeyURL != "" {
			fmt.Printf("获取 API Key: %s\n", p.KeyURL)
		}
		apiKey, err := ui.PromptAPIKey("请输入 API Key")
		if err != nil {
			return
		}
		p.APIKey = apiKey
		if err := config.AddProvider(p); err != nil {
			fmt.Fprintf(os.Stderr, "%s 保存配置失败: %v\n", red("错误:"), err)
			return
		}
		fmt.Printf("✓ 已配置 %s\n", p.DisplayName)
	case 2: // default
		if err := config.SetDefault(selectedName); err != nil {
			fmt.Fprintf(os.Stderr, "%s 设置默认失败: %v\n", red("错误:"), err)
			return
		}
		fmt.Printf("✓ 已设置 %s 为默认供应商\n", selectedName)
	case 3: // test
		testCmd.Run(nil, []string{selectedName})
	}
}

func init() {
	listCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "交互模式")
	rootCmd.AddCommand(listCmd)
}
