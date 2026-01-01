package cmd

import (
	"fmt"
	"os"

	"ccm/internal/config"
	"ccm/internal/provider"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	apiKey  string
	baseURL string
	model   string
)

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "添加或配置供应商",
	Long: `添加或配置供应商

预置供应商只需提供 --key:
  ccm add doubao --key "sk-xxx"

自定义供应商需要完整配置:
  ccm add custom --key "xxx" --url "https://..." --model "xxx"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		if apiKey == "" {
			fmt.Fprintln(os.Stderr, red("错误: 必须提供 --key 参数"))
			os.Exit(1)
		}

		var p provider.Provider

		// 检查是否是预置供应商
		if preset, ok := provider.Presets[name]; ok {
			p = preset
			p.APIKey = apiKey
			// 允许覆盖预置的 URL 和模型
			if baseURL != "" {
				p.BaseURL = baseURL
			}
			if model != "" {
				p.Model = model
			}
		} else {
			// 自定义供应商
			if baseURL == "" || model == "" {
				fmt.Fprintln(os.Stderr, red("错误: 自定义供应商必须提供 --url 和 --model 参数"))
				os.Exit(1)
			}
			p = provider.Provider{
				Name:        name,
				DisplayName: name,
				APIKey:      apiKey,
				BaseURL:     baseURL,
				Model:       model,
			}
		}

		if err := config.AddProvider(p); err != nil {
			fmt.Fprintf(os.Stderr, "%s 保存配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		fmt.Printf("%s 已配置供应商: %s\n", green("✓"), p.DisplayName)
		fmt.Printf("  API URL: %s\n", p.BaseURL)
		fmt.Printf("  模型: %s\n", p.Model)
		fmt.Println()
		fmt.Printf("使用 'ccm run %s' 启动 Claude Code\n", name)
	},
}

func init() {
	addCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API 密钥 (必填)")
	addCmd.Flags().StringVarP(&baseURL, "url", "u", "", "API URL (自定义供应商必填)")
	addCmd.Flags().StringVarP(&model, "model", "m", "", "模型名称 (自定义供应商必填)")
	rootCmd.AddCommand(addCmd)
}
