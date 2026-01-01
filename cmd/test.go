package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"ccm/internal/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test <name>",
	Short: "测试供应商 API 连接",
	Long: `测试供应商 API 连接是否正常

示例:
  ccm test doubao     测试豆包连接
  ccm test deepseek   测试 DeepSeek 连接`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
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
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 未配置\n", red("错误:"), name)
			fmt.Fprintf(os.Stderr, "请先配置: ccm add %s --key \"你的API密钥\"\n", name)
			os.Exit(1)
		}

		// 获取 API Key (支持环境变量)
		apiKey := p.APIKey
		envKey := config.GetEnvAPIKey(name)
		if envKey != "" {
			apiKey = envKey
		}

		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 未设置 API Key\n", red("错误:"), name)
			os.Exit(1)
		}

		fmt.Printf("测试供应商: %s (%s)\n", p.DisplayName, p.BaseURL)
		fmt.Println("正在连接...")

		// 创建 HTTP 请求测试连接
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		req, err := http.NewRequest("GET", p.BaseURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s 创建请求失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)

		start := time.Now()
		resp, err := client.Do(req)
		latency := time.Since(start)

		if err != nil {
			fmt.Printf("%s 连接失败\n", red("✗"))
			fmt.Printf("错误: %v\n", err)
			fmt.Printf("\n可能的原因:\n")
			fmt.Printf("  - API URL 错误: %s\n", p.BaseURL)
			fmt.Printf("  - API Key 无效或已过期\n")
			fmt.Printf("  - 网络连接问题\n")
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 || resp.StatusCode == 401 || resp.StatusCode == 403 {
			fmt.Printf("%s 连接成功!\n", green("✓"))
			fmt.Printf("  延迟: %v\n", latency)
			fmt.Printf("  状态码: %d\n", resp.StatusCode)

			if resp.StatusCode == 401 || resp.StatusCode == 403 {
				fmt.Printf("\n%s API Key 可能无效 (状态码 %d)\n", yellow("⚠️"), resp.StatusCode)
				fmt.Printf("建议: ccm edit %s --key \"新的API密钥\"\n", name)
			}
		} else {
			fmt.Printf("%s 连接异常\n", yellow("⚠️"))
			fmt.Printf("  状态码: %d\n", resp.StatusCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
