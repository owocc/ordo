package tui

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/owocc/ordo/internal/agent"
	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/opener"
	"github.com/owocc/ordo/internal/store"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		switch m.view {
		case viewProjects, viewIDEs:
			return m.updateList(msg)
	case viewAgentProjects:
		return m.updateAgentList(msg)
	case viewAddProject:
		return m.updateAddProject(msg)
	case viewAddIDE:
		return m.updateAddIDE(msg)
		case viewConfirmDelete:
			return m.updateConfirm(msg)
		}
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := len(m.projects)
	if m.view == viewIDEs {
		items = len(m.ides)
	}
	m.err, m.status = "", ""
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < items-1 {
			m.cursor++
		}
	case "tab":
		switch m.view {
		case viewProjects:
			m.view = viewIDEs
		case viewIDEs:
			m.view = viewAgentProjects
		default:
			m.view = viewProjects
		}
		m.cursor = 0
	case "r":
		m.projects = store.FindAllProjects()
		m.ides = store.FindAllIDEs()
		m.agentProjects, _ = agent.Discover()
		m.status = "已刷新"
	case "a":
		for i := range m.inputs {
			m.inputs[i].Reset()
		}
		m.inputs[0].Placeholder = "项目名称"
		m.inputs[1].Placeholder = "项目目录"
		m.inputs[2].Placeholder = ""
		m.view = viewAddProject
		m.focus = 0
		m.addProjIDE = -1
		m.inputs[0].Focus()
	case "A":
		for i := range m.inputs {
			m.inputs[i].Reset()
		}
		m.inputs[0].Placeholder = "IDE 名称"
		m.inputs[1].Placeholder = "命令或路径（如 code、cursor、/usr/bin/code）"
		m.inputs[2].Placeholder = "可选描述"
		m.view = viewAddIDE
		m.focus = 0
		m.inputs[0].Focus()
	case "d":
		if items == 0 {
			return m, nil
		}
		if m.view == viewIDEs {
			m.prevView = viewIDEs
			m.deleteIsIDE = true
			m.deleteTarget = m.ides[m.cursor].Name
		} else {
			m.prevView = viewProjects
			m.deleteIsIDE = false
			m.deleteTarget = m.projects[m.cursor].Name
		}
		m.view = viewConfirmDelete
	case "enter":
		if m.view == viewProjects && len(m.projects) > 0 {
			p := m.projects[m.cursor]
			ide := store.FindOneIDE(p.RelationIDEID)
			if ide == nil {
				m.err = "该项目未配置 IDE"
				return m, nil
			}
			command := opener.Run(opener.BuildCommand(p, ide))
			if err := command.Start(); err != nil {
				m.err = "打开失败: " + err.Error()
				return m, nil
			}
			m.lastOpened = p.Name
			m.status = "已打开: " + p.Name
		}
	}
	return m, nil
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if m.deleteIsIDE {
			if _, err := store.DeleteIDE(m.deleteTarget); err != nil {
				m.err = err.Error()
			}
		} else {
			if _, err := store.DeleteProject(m.deleteTarget); err != nil {
				m.err = err.Error()
			}
		}
		m.projects = store.FindAllProjects()
		m.ides = store.FindAllIDEs()
		m.view = m.prevView
		m.status = "已删除: " + m.deleteTarget
	case "n", "N", "esc":
		m.view = m.prevView
		m.deleteTarget = ""
	}
	return m, nil
}

func (m Model) updateAddProject(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.view = viewProjects
		m.cursor = 0
		m.projects = store.FindAllProjects()
		return m, nil
	}
	if m.focus < 2 {
		switch msg.String() {
		case "enter", "tab":
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus < 2 {
				m.inputs[m.focus].Focus()
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}
	// IDE 选择阶段
	switch msg.String() {
	case "up", "k":
		if m.addProjIDE > -1 {
			m.addProjIDE--
		}
	case "down", "j":
		if m.addProjIDE < len(m.ides)-1 {
			m.addProjIDE++
		}
	case "enter":
		name := strings.TrimSpace(m.inputs[0].Value())
		dir := strings.TrimSpace(m.inputs[1].Value())
		if name == "" || dir == "" {
			m.err = "名称和目录必填"
			return m, nil
		}
		if ve := model.ValidateName(name); ve != nil {
			m.err = ve.Error()
			return m, nil
		}
		if !filepath.IsAbs(dir) {
			if abs, err := filepath.Abs(dir); err == nil {
				dir = abs
			}
		}
		var ideID string
		if m.addProjIDE >= 0 && m.addProjIDE < len(m.ides) {
			ideID = m.ides[m.addProjIDE].ID
		}
		if _, err := store.AddProject(model.Project{
			Name:          name,
			Dir:           dir,
			RelationIDEID: ideID,
		}); err != nil {
			m.err = err.Error()
			return m, nil
		}
		m.projects = store.FindAllProjects()
		m.view = viewProjects
		m.cursor = len(m.projects) - 1
		m.status = "已添加项目: " + name
	}
	return m, nil
}

func (m Model) updateAddIDE(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.view = viewIDEs
		m.cursor = 0
		m.ides = store.FindAllIDEs()
		return m, nil
	}
	switch msg.String() {
	case "enter", "tab":
		m.inputs[m.focus].Blur()
		m.focus++
		if m.focus < 3 {
			m.inputs[m.focus].Focus()
			return m, nil
		}
		// 提交
		name := strings.TrimSpace(m.inputs[0].Value())
		path := strings.TrimSpace(m.inputs[1].Value())
		desc := strings.TrimSpace(m.inputs[2].Value())
		if name == "" || path == "" {
			m.err = "名称和路径必填"
			m.focus = 0
			m.inputs[0].Focus()
			return m, nil
		}
		if ve := model.ValidateName(name); ve != nil {
			m.err = ve.Error()
			m.focus = 0
			m.inputs[0].Focus()
			return m, nil
		}
		if _, err := store.AddIDE(model.IDE{Name: name, Path: path, Desc: desc}); err != nil {
			m.err = err.Error()
			m.focus = 0
			m.inputs[0].Focus()
			return m, nil
		}
		m.ides = store.FindAllIDEs()
		m.view = viewIDEs
		m.cursor = len(m.ides) - 1
		m.status = "已添加 IDE: " + name
		return m, nil
	}
	if m.focus < 3 {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) updateAgentList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.err, m.status = "", ""
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.agentProjects)-1 {
			m.cursor++
		}
	case "tab":
		m.view = viewProjects
		m.cursor = 0
	case "r":
		m.agentProjects, _ = agent.Discover()
		m.status = "已刷新"
	case "enter":
		if len(m.agentProjects) == 0 {
			return m, nil
		}
		p := m.agentProjects[m.cursor]
		ides := store.FindAllIDEs()
		if len(ides) == 0 {
			m.err = "未注册 IDE，请按 A 或 Tab 切换到 IDE 视图添加"
			return m, nil
		}
		// 简单的打开策略：使用第一个已注册的 IDE
		modelP := agent.ToProject(&p)
		ide := &ides[0]
		command := opener.Run(opener.BuildCommand(modelP, ide))
		if err := command.Start(); err != nil {
			m.err = "打开失败: " + err.Error()
			return m, nil
		}
		m.status = "已打开: " + p.Path
	}
	return m, nil
}
