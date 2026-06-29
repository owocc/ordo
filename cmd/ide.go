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
		name, path := "", ""
		if len(args) > 0 {
			name = args[0]
		}
		if len(args) > 1 {
			path = args[1]
		}
		if name == "" {
			name = promptText("请输入 IDE 名称（例如: VSCode）：")
			if name == "" {
				return fail("IDE 名称是必填的", nil)
			}
		}
		if path == "" {
			path = promptText("请输入 IDE 路径或命令（例如: /usr/bin/code 或 idea）：")
			if path == "" {
				return fail("IDE 路径或命令是必填的", nil)
			}
		}
		desc := promptText("可选 - IDE 描述：")
		if _, err := store.AddIDE(model.IDE{Name: name, Path: path, Desc: desc}); err != nil {
			return fail("添加 IDE 失败", err)
		}
		ui.Success("IDE '" + name + "' 已成功添加！")
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
