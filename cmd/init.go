package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"ccm/internal/config"
	"ccm/internal/provider"
	"ccm/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "首次使用引导",
	Long: `首次使用引导

帮助您完成初始配置:
1. 检测环境 (npm, claude)
2. 选择并配置供应商
3. 测试连接`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Println("欢迎使用 CCM (Claude Code Manager)!")
		fmt.Println()

		// 1. 检测环境
		fmt.Println("检测环境...")

		// 检查 npm
		hasNPM := func() bool {
			_, err := exec.LookPath("npm")
			return err == nil
		}()

		if hasNPM {
			fmt.Printf("  %s npm 已安装\n", green("✓"))
		} else {
			fmt.Printf("  %s npm 未安装\n", yellow("✗"))
			fmt.Println("  请先安装 Node.js: https://nodejs.org/")
		}

		// 检查 claude
		hasClaude := func() bool {
			_, err := exec.LookPath("claude")
			return err == nil
		}()

		home, _ := os.UserHomeDir()
		localClaude := filepath.Join(home, "claude-model", "node_modules", ".bin", "claude")
		if _, err := os.Stat(localClaude); err == nil {
			hasClaude = true
		}

		if hasClaude {
			fmt.Printf("  %s Claude Code 已安装\n", green("✓"))
		} else {
			fmt.Printf("  %s Claude Code 未安装\n", yellow("✗"))
			fmt.Println("  运行 ccm run <provider> 时会自动提示安装")
		}

		fmt.Println()

		// 2. 构建供应商列表
		cfg, _ := config.Load()
		items := ui.BuildProviderItems(cfg, true)

		// 3. 交互式选择
		if !ui.PromptConfirm("是否立即配置一个供应商") {
			fmt.Println()
			fmt.Println("快速开始:")
			fmt.Printf("  %s                # 查看所有供应商\n", cyan("ccm list"))
			fmt.Printf("  %s  # 配置豆包\n", cyan("ccm add doubao --key \"xxx\""))
			fmt.Printf("  %s          # 启动 Claude\n", cyan("ccm run doubao"))
			return
		}

		fmt.Println()

		// 使用箭头键选择供应商
		selectedName, err := ui.SelectProvider(items, "选择要配置的供应商")
		if err != nil {
			fmt.Println("取消选择")
			return
		}

		p := provider.Presets[selectedName]

		fmt.Printf("\n您选择了: %s (%s)\n", cyan(selectedName), p.DisplayName)
		fmt.Println("API URL:", p.BaseURL)
		fmt.Println("默认模型:", p.Model)
		if p.KeyURL != "" {
			fmt.Printf("获取 API Key: %s\n", cyan(p.KeyURL))
		}
		fmt.Println()

		// 使用掩码输入 API Key
		apiKey, err := ui.PromptAPIKey("请输入 API Key")
		if err != nil {
			fmt.Println("取消输入")
			return
		}

		// 保存配置
		p.APIKey = apiKey
		if err := config.AddProvider(p); err != nil {
			fmt.Printf("保存配置失败: %v\n", err)
			return
		}

		fmt.Printf("\n%s 已配置 %s!\n", green("✓"), p.DisplayName)
		fmt.Println()
		fmt.Println("下一步:")
		fmt.Printf("  %s     # 测试连接\n", cyan(fmt.Sprintf("ccm test %s", selectedName)))
		fmt.Printf("  %s      # 启动 Claude Code\n", cyan(fmt.Sprintf("ccm run %s", selectedName)))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
