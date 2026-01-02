package cmd

import (
	"fmt"
	"runtime"

	"ccm/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 CCM 的版本信息，包括版本号、Commit ID、构建时间和 Go 版本`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CCM (Claude Code Manager)")
		fmt.Println("=========================")
		fmt.Printf("Version:   %s\n", version.Version)
		fmt.Printf("Commit:    %s\n", version.Commit)
		fmt.Printf("Build Date: %s\n", version.Date)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
