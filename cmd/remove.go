package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"ccm/internal/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var forceRemove bool

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "删除已配置的供应商",
	Long: `删除已配置的供应商

示例:
  ccm remove doubao     删除豆包配置
  ccm remove deepseek   删除 DeepSeek 配置`,
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
			fmt.Fprintf(os.Stderr, "%s 供应商 '%s' 不存在\n", red("错误:"), name)
			os.Exit(1)
		}

		// 确认删除
		if !forceRemove {
			fmt.Printf("%s 确认删除供应商 '%s' (%s)? [y/N]: ", yellow("⚠️"), name, p.DisplayName)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if input != "y\n" && input != "Y\n" {
				fmt.Println("已取消")
				os.Exit(0)
			}
		}

		// 删除供应商配置
		delete(cfg.Providers, name)

		// 保存配置
		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "%s 保存配置失败: %v\n", red("错误:"), err)
			os.Exit(1)
		}

		// 删除对应的配置目录
		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, "claude-model", "configs", ".claude-"+name)
		if _, err := os.Stat(configDir); err == nil {
			os.RemoveAll(configDir)
		}

		fmt.Printf("%s 已删除供应商: %s\n", green("✓"), p.DisplayName)
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "强制删除，不询问确认")
	rootCmd.AddCommand(removeCmd)
}
