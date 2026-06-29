package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/store"
)

func (m Model) renderList(projectsMode bool) string {
	var b strings.Builder
	if projectsMode {
		b.WriteString(titleProj.Render(fmt.Sprintf("Ordo · 项目 (%d)", len(m.projects))))
	} else {
		b.WriteString(titleIDE.Render(fmt.Sprintf("Ordo · IDE (%d)", len(m.ides))))
	}
	b.WriteString("\n\n")

	if projectsMode {
		if len(m.projects) == 0 {
			b.WriteString(hintStyle.Render("  没有项目。按 a 添加一个项目。"))
			b.WriteString("\n")
		} else {
			for i, p := range m.projects {
				line := fmt.Sprintf("  %-20s %s", p.DisplayOrName(), truncate(p.Dir, 50))
				ideName := ideNameFor(p)
				if ideName != "" {
					line += "  [" + ideName + "]"
				}
				if i == m.cursor {
					b.WriteString(selStyle.Render("> " + line))
				} else {
					b.WriteString(normalStyle.Render("  " + line))
				}
				b.WriteString("\n")
			}
		}
		b.WriteString("\n" + hintStyle.Render(projectsHelp))
	} else {
		if len(m.ides) == 0 {
			b.WriteString(hintStyle.Render("  没有 IDE。按 A 添加一个 IDE。"))
			b.WriteString("\n")
		} else {
			for i, ide := range m.ides {
				line := fmt.Sprintf("  %-20s %s", ide.DisplayOrName(), truncate(ide.Path, 50))
				if i == m.cursor {
					b.WriteString(selStyle.Render("> " + line))
				} else {
					b.WriteString(normalStyle.Render("  " + line))
				}
				b.WriteString("\n")
			}
		}
		b.WriteString("\n" + hintStyle.Render(idesHelp))
	}

	if m.status != "" {
		b.WriteString("\n" + okStyle.Render(m.status))
	}
	if m.err != "" {
		b.WriteString("\n" + errStyle.Render(m.err))
	}
	if m.lastOpened != "" {
		b.WriteString("\n" + okStyle.Render("已打开: "+m.lastOpened))
	}
	return b.String()
}

func (m Model) renderAddProject() string {
	var b strings.Builder
	b.WriteString(titleProj.Render("Ordo · 添加项目"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  名称:  %s\n", m.inputs[0].View()))
	b.WriteString(fmt.Sprintf("  目录:  %s\n", m.inputs[1].View()))
	// IDE 选项
	ides := m.ides
	line := "  IDE:   "
	if m.focus == 2 {
		line = "  IDE:  >"
	}
	ideLabel := "无 IDE"
	if m.addProjIDE >= 0 && m.addProjIDE < len(ides) {
		ideLabel = ides[m.addProjIDE].Name
	}
	line += ideLabel + "  (↑/↓ 选择, Enter 确认并添加)"
	b.WriteString(line + "\n\n")
	b.WriteString(hintStyle.Render(formHelp))
	if m.err != "" {
		b.WriteString("\n" + errStyle.Render(m.err))
	}
	return b.String()
}

func (m Model) renderAddIDE() string {
	var b strings.Builder
	b.WriteString(titleIDE.Render("Ordo · 添加 IDE"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  名称:    %s\n", m.inputs[0].View()))
	b.WriteString(fmt.Sprintf("  命令/路径: %s\n", m.inputs[1].View()))
	b.WriteString(fmt.Sprintf("  描述:    %s\n", m.inputs[2].View()))
	b.WriteString("\n" + hintStyle.Render("Enter 提交    Esc 取消"))
	if m.err != "" {
		b.WriteString("\n" + errStyle.Render(m.err))
	}
	return b.String()
}

func (m Model) renderConfirmDelete() string {
	kind := "项目"
	if m.deleteIsIDE {
		kind = "IDE"
	}
	var b strings.Builder
	b.WriteString(titleProj.Render("Ordo · 确认删除"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  确认删除 %s: %s ?\n\n", kind, m.deleteTarget))
	b.WriteString(hintStyle.Render("  y 确认删除    n/Esc 取消"))
	return b.String()
}

func (m Model) renderAgentList() string {
	var b strings.Builder
	b.WriteString(titleAgent.Render(fmt.Sprintf("Ordo · Agent 项目 (%d)", len(m.agentProjects))))
	b.WriteString("\n\n")
	if len(m.agentProjects) == 0 {
		b.WriteString(hintStyle.Render("  没有发现 agent 项目。按 r 刷新。"))
		b.WriteString("\n")
	} else {
		for i, p := range m.agentProjects {
			sources := strings.Join(p.Sources, "+")
			age := relativeTimeAgent(p.LastActive)
			line := fmt.Sprintf("  %-55s [%s] %s %s",
				truncate(p.Path, 55),
				sourceStyle.Render(sources),
				sessionStyle.Render(fmt.Sprintf("%3d", p.SessionCount)),
				age,
			)
			if i == m.cursor {
				b.WriteString(selStyle.Render("> " + line))
			} else {
				b.WriteString(normalStyle.Render("  " + line))
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("\n" + hintStyle.Render(agentHelp))
	if m.status != "" {
		b.WriteString("\n" + okStyle.Render(m.status))
	}
	if m.err != "" {
		b.WriteString("\n" + errStyle.Render(m.err))
	}
	return b.String()
}

func relativeTimeAgent(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "刚刚"
	case d < time.Hour:
		return fmt.Sprintf("%dm前", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh前", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd前", int(d.Hours()/24))
	default:
		return t.Format("01-02")
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}

func ideNameFor(p model.Project) string {
	if p.RelationIDEID == "" {
		return ""
	}
	ide := store.FindOneIDE(p.RelationIDEID)
	if ide == nil {
		return ""
	}
	return ide.DisplayOrName()
}
