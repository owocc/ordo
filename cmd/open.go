package cmd

import (
	"fmt"

	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/opener"
	"github.com/owocc/ordo/internal/store"
	"github.com/owocc/ordo/internal/ui"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [idOrName...]",
	Short: "在关联的 IDE 中打开指定项目",
	Long:  "通过项目名称或 ID 打开项目。可一次传入多个项目。未提供参数时进入选择模式。",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		projects := store.FindAllProjects()
		if len(projects) == 0 {
			ui.Warn("当前没有已管理的项目。请先使用 'ordo add' 添加。")
			return nil
		}

		var selected []model.Project
		if len(args) == 0 {
			// 无参数时让用户从列表选择
			selected = promptSelectProjects(projects)
		} else {
			for _, idOrName := range args {
				p := store.FindOneProject(idOrName)
				if p == nil {
					ui.Warn("未找到项目: " + idOrName)
					continue
				}
				selected = append(selected, *p)
			}
		}
		if len(selected) == 0 {
			ui.Warn("没有匹配的项目。")
			return nil
		}

		ides := store.FindAllIDEs()
		if len(ides) == 0 {
			return fail("未找到 IDE。请在继续之前配置至少一个 IDE。", nil)
		}
		ideMap := make(map[string]*model.IDE, len(ides))
		for i := range ides {
			ideMap[ides[i].ID] = &ides[i]
		}

		for _, p := range selected {
			ide := ideMap[p.RelationIDEID]
			if ide == nil {
				// 无关联 IDE，弹选择
				choice := promptSelectIDE("🤯 项目 '"+p.Name+"' 没有配置默认 IDE，请选择一个：", "")
				if choice == "" {
					ui.Warn("跳过项目 " + p.Name)
					continue
				}
				ide = store.FindOneIDE(choice)
			}
			command := opener.BuildCommand(p, ide)
			ui.Info("执行命令: " + command)
			c := opener.Run(command)
			if err := c.Start(); err != nil {
				return fail(fmt.Sprintf("打开项目 '%s' 失败", p.Name), err)
			}
			ui.Success("已启动项目 '" + p.Name + "'")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
