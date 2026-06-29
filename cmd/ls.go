package cmd

import (
	"github.com/owocc/ordo/internal/ui"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有已管理的项目",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintProjects()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
