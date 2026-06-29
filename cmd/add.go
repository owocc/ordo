package cmd

import (
	"path/filepath"

	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/store"
	"github.com/owocc/ordo/internal/ui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name] [dir] [ide]",
	Short: "将一个项目添加到管理器",
	Long: "将一个项目添加到管理器，指定项目名称、目录，并可选地关联一个 IDE。\n" +
		"未提供的参数会进入交互式输入。",
	Args: cobra.MaximumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, dir, ideArg := "", "", ""
		if len(args) > 0 {
			name = args[0]
		}
		if len(args) > 1 {
			dir = args[1]
		}
		if len(args) > 2 {
			ideArg = args[2]
		}

		if name == "" {
			name = promptText("请输入项目名称：")
			if name == "" {
				return fail("项目名称是必填项", nil)
			}
		}
		if dir == "" {
			dir = promptText("请输入项目目录：")
			if dir == "" {
				return fail("项目目录是必填项", nil)
			}
		}

		// 解析相对路径为绝对路径
		if !filepath.IsAbs(dir) {
			if abs, err := filepath.Abs(dir); err == nil {
				dir = abs
			}
		}

		// IDE 选择
		var ideID string
		if ideArg != "" {
			ide := store.FindOneIDE(ideArg)
			if ide == nil {
				return fail("没有找到标识符为 \""+ideArg+"\" 的 IDE", nil)
			}
			ideID = ide.ID
		} else {
			ideID = promptSelectIDE("请选择一个 IDE（或跳过）：", "")
		}

		if _, err := store.AddProject(model.Project{
			Name:          name,
			Dir:           dir,
			RelationIDEID: ideID,
		}); err != nil {
			return fail("添加项目失败", err)
		}
		ui.Success("项目 '" + name + "' 已成功添加！")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
