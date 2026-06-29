package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/owocc/ordo/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "进入交互式 TUI 模式，浏览并管理项目和 IDE",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.New(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fail("启动 TUI 失败", fmt.Errorf("%w", err))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
