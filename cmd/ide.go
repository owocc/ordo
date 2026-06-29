package cmd

import (
	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/store"
	"github.com/owocc/ordo/internal/ui"
	"github.com/spf13/cobra"
)

var ideCmd = &cobra.Command{
	Use:   "ide",
	Short: "管理已注册的 IDE",
}

var ideAddCmd = &cobra.Command{
	Use:   "add [name] [pathOrCommand]",
	Short: "将一个新的 IDE 添加到管理器",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, path, displayName string
		if len(args) > 0 {
			name = args[0]
		}
		if len(args) > 1 {
			path = args[1]
		}
		if name == "" {
			var ve error
			name, displayName, ve = promptNameAndDisplay("IDE")
			if ve != nil {
				return fail(ve.Error(), nil)
			}
		} else {
			if ve := model.ValidateName(name); ve != nil {
				return fail("无效的 IDE 名称", ve)
			}
			displayName = promptText("显示名称（可选，支持空格，直接回车使用 '" + name + "'）：")
		}
		if displayName == "" {
			displayName = name
		}
		if path == "" {
			path = promptText("请输入 IDE 路径或命令（例如: code、cursor、或 /usr/bin/code）：")
			if path == "" {
				return fail("IDE 路径或命令是必填的", nil)
			}
		}
		desc := promptText("可选 - IDE 描述：")
		if _, err := store.AddIDE(model.IDE{
			Name:        name,
			DisplayName: displayName,
			Path:        path,
			Desc:        desc,
		}); err != nil {
			return fail("添加 IDE 失败", err)
		}
		ui.Success("IDE '" + displayName + "' 已成功添加！")
		return nil
	},
}

var ideLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有已注册的 IDE",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintIDEs()
		return nil
	},
}

var ideRmCmd = &cobra.Command{
	Use:   "rm [idOrName]",
	Short: "通过 ID 或名称删除一个 IDE",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mark := ""
		if len(args) > 0 {
			mark = args[0]
		} else {
			mark = promptText("请输入要删除的 IDE 的名称或 ID：")
		}
		if mark == "" {
			return fail("IDE 标识符是必填项", nil)
		}
		if store.FindOneIDE(mark) == nil {
			return fail("没有找到标识符为 \""+mark+"\" 的 IDE", nil)
		}
		if _, err := store.DeleteIDE(mark); err != nil {
			return fail("删除 IDE 失败", err)
		}
		ui.Success("已删除 IDE: " + mark)
		return nil
	},
}

func init() {
	ideCmd.AddCommand(ideAddCmd, ideLsCmd, ideRmCmd)
	rootCmd.AddCommand(ideCmd)
}
