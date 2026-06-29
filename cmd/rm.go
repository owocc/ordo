package cmd

import (
	"github.com/owocc/ordo/internal/store"
	"github.com/owocc/ordo/internal/ui"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [idOrName]",
	Short: "从管理器中删除一个项目",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mark := ""
		if len(args) > 0 {
			mark = args[0]
		} else {
			mark = promptText("请输入要删除的项目名称或 ID：")
		}
		if mark == "" {
			return fail("项目标识符是必填项", nil)
		}
		removed, err := store.DeleteProject(mark)
		if err != nil {
			return fail("删除项目失败", err)
		}
		if removed == 0 {
			ui.Warn("未找到标识符为 \"" + mark + "\" 的项目")
			return nil
		}
		ui.Success("已删除项目: " + mark)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
