package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/store"
)

// 视图状态
type view int

const (
	viewProjects view = iota
	viewIDEs
	viewAddProject
	viewAddIDE
	viewConfirmDelete
)

type focusField int

var (
	titleProj   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Padding(0, 1)
	titleIDE    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).Padding(0, 1)
	selStyle    = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("15")).Bold(true)
	normalStyle = lipgloss.NewStyle()
	hintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	errStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
)

const projectsHelp = `↑/↓ 选择    Enter 打开    d 删除    a 添加项目    A 添加 IDE    Tab 切换视图    q 退出`
const idesHelp = `↑/↓ 选择    d 删除    A 添加 IDE    a 添加项目    Tab 切换视图    q 退出`
const formHelp = `Tab/Enter 下一步    Esc 取消`

// Model 是 TUI 的顶层状态。
type Model struct {
	view     view
	prevView view
	cursor   int
	projects []model.Project
	ides     []model.IDE
	width    int
	height   int
	err      string
	status   string
	// 表单
	inputs []textinput.Model
	focus  focusField
	// 添加项目时选定的 ide index（-1 = 无）
	addProjIDE int
	// 删除确认
	deleteTarget string
	deleteIsIDE  bool
	// 上次打开的项目名
	lastOpened string
}

// New 创建初始 Model。
func New() Model {
	ti := func(placeholder string) textinput.Model {
		t := textinput.New()
		t.Placeholder = placeholder
		t.CharLimit = 256
		t.Width = 40
		return t
	}
	return Model{
		view:       viewProjects,
		projects:   store.FindAllProjects(),
		ides:       store.FindAllIDEs(),
		inputs:     []textinput.Model{ti("项目名称"), ti("项目目录"), ti("")},
		addProjIDE: -1,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) View() string {
	switch m.view {
	case viewProjects:
		return m.renderList(true)
	case viewIDEs:
		return m.renderList(false)
	case viewAddProject:
		return m.renderAddProject()
	case viewAddIDE:
		return m.renderAddIDE()
	case viewConfirmDelete:
		return m.renderConfirmDelete()
	}
	return ""
}
