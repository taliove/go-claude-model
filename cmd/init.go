package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ccm/internal/config"
	"ccm/internal/provider"

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

		fmt.Println("👋 欢迎使用 CCM (Claude Code Manager)!")
		fmt.Println()

		// 1. 检测环境
		fmt.Println("🔍 检测环境...")

		// 检查 npm
		hasNPM := func() bool {
			_, err := exec.LookPath("npm")
			return err == nil
		}()

		if hasNPM {
			fmt.Printf("  %s npm 已安装\n", green("✓"))
		} else {
			fmt.Printf("  %s npm 未安装\n", yellow("✗"))
			fmt.Println("  💡 请先安装 Node.js: https://nodejs.org/")
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
			fmt.Println("  💡 运行 ccm run <provider> 时会自动提示安装")
		}

		fmt.Println()

		// 2. 显示可用供应商
		fmt.Println("📦 可用的供应商:")
		fmt.Println()

		presets := []string{}
		for name := range provider.Presets {
			presets = append(presets, name)
		}

		for i, name := range presets {
			p := provider.Presets[name]
			cfg, _ := config.Load()
			if _, ok := cfg.Providers[name]; ok {
				fmt.Printf("  %d. %s (%s) %s\n", i+1, cyan(name), p.DisplayName, green("✓ 已配置"))
			} else {
				fmt.Printf("  %d. %s (%s)\n", i+1, cyan(name), p.DisplayName)
			}
		}

		fmt.Println()
		fmt.Println("💡 使用 'ccm add <name> --key \"your-api-key\"' 配置供应商")

		// 3. 交互式选择
		fmt.Println()
		fmt.Print("是否立即配置一个供应商? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if input != "y\n" && input != "Y\n" {
			fmt.Println()
			fmt.Println("🚀 快速开始:")
			fmt.Println("  ccm list                # 查看所有供应商")
			fmt.Println("  ccm add doubao --key \"xxx\"  # 配置豆包")
			fmt.Println("  ccm run doubao          # 启动 Claude")
			return
		}

		fmt.Print("\n请选择供应商编号: ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var selectedIdx int
		fmt.Sscanf(input, "%d", &selectedIdx)

		if selectedIdx < 1 || selectedIdx > len(presets) {
			fmt.Println("无效的选择")
			return
		}

		selectedName := presets[selectedIdx-1]
		p := provider.Presets[selectedName]

		fmt.Printf("\n您选择了: %s (%s)\n", selectedName, p.DisplayName)
		fmt.Println("API URL:", p.BaseURL)
		fmt.Println("默认模型:", p.Model)
		fmt.Println()

		fmt.Print("请输入 API Key: ")
		input, _ = reader.ReadString('\n')
		apiKey := strings.TrimSpace(input)

		if apiKey == "" {
			fmt.Println("API Key 不能为空")
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
		fmt.Println("📖 下一步:")
		fmt.Printf("  ccm test %s     # 测试连接\n", selectedName)
		fmt.Printf("  ccm run %s      # 启动 Claude Code\n", selectedName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
