package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/owocc/ordo/internal/agent"
	"github.com/owocc/ordo/internal/opener"
	"github.com/owocc/ordo/internal/store"
	"github.com/owocc/ordo/internal/ui"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "从 agent 应用（opencode / codex / claude）发现并打开项目",
}

var agentLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出从 agent 应用发现的项目目录",
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := agent.Discover()
		if err != nil {
			return fail("发现 agent 项目失败", err)
		}
		ui.PrintAgent(projects)
		return nil
	},
}

var agentOpenCmd = &cobra.Command{
	Use:   "open [idOrPath]",
	Short: "用 IDE 打开从 agent 发现的项目目录",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := agent.Discover()
		if err != nil {
			return fail("发现 agent 项目失败", err)
		}
		if len(projects) == 0 {
			ui.Info("没有发现 agent 项目。")
			return nil
		}

		var target *agent.Project
		if len(args) == 0 {
			// 交互选择
			target = promptSelectAgentProject(projects)
		} else {
			arg := args[0]
			// 先试序号
			if n, err := strconv.Atoi(arg); err == nil && n >= 1 && n <= len(projects) {
				target = &projects[n-1]
			} else {
				// 再按路径前缀匹配
				for _, p := range projects {
					if strings.Contains(p.Path, arg) {
						target = &p
						break
					}
				}
				if target == nil {
					return fail("未找到匹配的项目: "+arg, nil)
				}
			}
		}
		if target == nil {
			return nil
		}

		// 选 IDE
		ides := store.FindAllIDEs()
		if len(ides) == 0 {
			return fail("未找到 IDE。请先使用 'ordo ide add' 注册一个 IDE。", nil)
		}
		ideID := promptSelectIDE("选择 IDE 打开项目: "+target.Path, "")

		// 通过 opener 打开
		ide := store.FindOneIDE(ideID)
		if ide == nil {
			return fail("选择的 IDE 不存在", nil)
		}
		p := agent.ToProject(target)
		cmdLine := opener.BuildCommand(p, ide)
		ui.Info("执行命令: " + cmdLine)
		c := opener.Run(cmdLine)
		if err := c.Start(); err != nil {
			return fail("打开失败", err)
		}
		ui.Success("已启动项目: " + target.Path)
		return nil
	},
}

func init() {
	agentCmd.AddCommand(agentLsCmd, agentOpenCmd)
	rootCmd.AddCommand(agentCmd)
}

// promptSelectAgentProject 列出项目让用户通过#号选择。
func promptSelectAgentProject(projects []agent.Project) *agent.Project {
	fmt.Println("请选择要打开的项目：")
	for i, p := range projects {
		sources := strings.Join(p.Sources, "+")
		fmt.Printf("  %d. %s  [%s]  (%d sessions)\n", i+1, p.Path, sources, p.SessionCount)
	}
	for {
		s := promptText("输入序号：")
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > len(projects) {
			fmt.Println("无效输入，请重新选择。")
			continue
		}
		return &projects[n-1]
	}
}
