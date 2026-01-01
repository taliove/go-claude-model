package cmd

import (
	"fmt"
	"os"

	"ccm/internal/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	newAPIKey string
	newBaseURL string
	newModel string
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "更新供应商配置",
	Long: `更新供应商配置

示例:
  ccm edit doubao --key "new-key"          更新 API Key
  ccm edit doubao --url "https://..."      更新 API URL
  ccm edit doubao --model "xxx"            更新模型
  ccm edit doubao -k "xxx" -u "..." -m "..."  一次性更新多个`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 加载配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		// 检查供应商是否存在
		p, ok := cfg.Providers[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 不存在\n", red("错误:"), name)
			fmt.Fprintf(os.Stderr, "可用供应商: ccm list\n")
			os.Exit(1)
		}

		// 更新字段
		updated := false
		if newAPIKey != "" {
			p.APIKey = newAPIKey
			updated = true
		}
		if newBaseURL != "" {
			p.BaseURL = newBaseURL
			updated = true
		}
		if newModel != "" {
			p.Model = newModel
			updated = true
		}

		if !updated {
			fmt.Fprintf(os.Stderr, "%s 请指定要更新的字段 (--key, --url, --model)\n", red("错误:"))
			os.Exit(1)
		}

		// 保存配置
		cfg.Providers[name] = p
		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "%s 保存配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		fmt.Printf("%s 已更新供应商: %s\n", green("✓"), p.DisplayName)
		if newAPIKey != "" {
			fmt.Printf("  API Key:    ********\n")
		}
		if newBaseURL != "" {
			fmt.Printf("  API URL:    %s\n", p.BaseURL)
		}
		if newModel != "" {
			fmt.Printf("  模型:       %s\n", p.Model)
		}
		fmt.Println()
		fmt.Printf("使用 'ccm run %s' 测试配置\n", name)
	},
}

func init() {
	editCmd.Flags().StringVarP(&newAPIKey, "key", "k", "", "新的 API 密钥")
	editCmd.Flags().StringVarP(&newBaseURL, "url", "u", "", "新的 API URL")
	editCmd.Flags().StringVarP(&newModel, "model", "m", "", "新的模型名称")
	rootCmd.AddCommand(editCmd)
}
