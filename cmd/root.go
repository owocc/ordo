package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ordo",
	Short: "Ordo — 简洁高效的项目与 IDE 管理命令行工具",
	Long: "Ordo 是一个跨平台的项目管理 CLI：\n" +
		"  添加 / 打开 / 列出 / 删除项目，并关联到已注册的 IDE。\n" +
		"  支持交互式 TUI 模式（ordo tui）。",
	Version: "2.0.0",
}

// Execute 是入口。
func Execute() error {
	return rootCmd.Execute()
}

// fail 是命令 action 的统一错误处理：打印并返回 err 让 cobra 退出码非 0。
func fail(msg string, err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ %s: %v\n", msg, err)
		return fmt.Errorf("%s: %w", msg, err)
	}
	fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
	return fmt.Errorf("%s", msg)
}
