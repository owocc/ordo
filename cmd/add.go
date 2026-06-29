package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/store"
	"github.com/owocc/ordo/internal/ui"
	"github.com/spf13/cobra"
)

func promptNameAndDisplay(what string) (name, display string, err error) {
	for {
		n := promptText("请输入" + what + "名称（标识符，不含空格，如 my-project）：")
		if n == "" {
			return "", "", fmt.Errorf("%s名称是必填项", what)
		}
		if ve := model.ValidateName(n); ve != nil {
			ui.Error(ve.Error())
			continue
		}
		name = n
		break
	}
	d := promptText("显示名称（可选，支持空格，直接回车使用 '" + name + "'）：")
	display = d
	return
}

var addCmd = &cobra.Command{
	Use:   "add [name] [dir] [ide]",
	Short: "将一个项目添加到管理器",
	Long: "将一个项目添加到管理器，指定项目名称、目录，并可选地关联一个 IDE。\n" +
		"未提供的参数会进入交互式输入。",
	Args: cobra.MaximumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			name        string
			displayName string
			dir         string
			ideArg      string
		)
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
			var ve error
			name, displayName, ve = promptNameAndDisplay("项目")
			if ve != nil {
				return fail(ve.Error(), nil)
			}
		} else {
			if ve := model.ValidateName(name); ve != nil {
				return fail("无效的项目名称", ve)
			}
			displayName = promptText("显示名称（可选，支持空格，直接回车使用 '" + name + "'）：")
		}
		if displayName == "" {
			displayName = name
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
			DisplayName:   displayName,
			Dir:           dir,
			RelationIDEID: ideID,
		}); err != nil {
			return fail("添加项目失败", err)
		}
		ui.Success("项目 '" + displayName + "' 已成功添加！")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
