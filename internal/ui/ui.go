package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/store"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
)

func Success(msg string) { fmt.Println(successStyle.Render("✓ " + msg)) }
func Info(msg string)    { fmt.Println(infoStyle.Render("• " + msg)) }
func Warn(msg string)    { fmt.Println(warnStyle.Render("! " + msg)) }
func Error(msg string)   { fmt.Fprintln(os.Stderr, errorStyle.Render("✗ "+msg)) }

// renderTable 渲染无边框、列间用空格分隔的表格（兼容旧 blankBorder 风格）。
func renderTable(headers []string, rows [][]string) string {
	cols := len(headers)
	widths := make([]int, cols)
	for i, h := range headers {
		widths[i] = lipgloss.Width(h)
	}
	for _, r := range rows {
		for i := 0; i < cols && i < len(r); i++ {
			if w := lipgloss.Width(r[i]); w > widths[i] {
				widths[i] = w
			}
		}
	}
	pad := func(s string, w int) string {
		return s + strings.Repeat(" ", w-lipgloss.Width(s))
	}
	var b strings.Builder
	hr := headerStyle.Render(strings.Join(func() []string {
		out := make([]string, cols)
		for i, h := range headers {
			out[i] = pad(h, widths[i])
		}
		return out
	}(), "   "))
	b.WriteString(hr)
	b.WriteString("\n")
	for _, r := range rows {
		cells := make([]string, cols)
		for i := 0; i < cols; i++ {
			val := ""
			if i < len(r) {
				val = r[i]
			}
			cells[i] = pad(val, widths[i])
		}
		b.WriteString(strings.Join(cells, "   "))
		b.WriteString("\n")
	}
	return b.String()
}

// PrintProjects 打印项目列表。
func PrintProjects() {
	projects := store.FindAllProjects()
	ides := store.FindAllIDEs()
	ideMap := make(map[string]*model.IDE, len(ides))
	for i := range ides {
		ideMap[ides[i].ID] = &ides[i]
	}
	if len(projects) == 0 {
		Info("没有找到任何项目。请使用 'ordo add' 添加一个项目。")
		return
	}
	rows := make([][]string, 0, len(projects))
	for _, p := range projects {
		ideName := "None"
		if p.RelationIDEID != "" {
			if ide, ok := ideMap[p.RelationIDEID]; ok && ide != nil {
				ideName = ide.Name
			}
		}
		rows = append(rows, []string{p.ID, p.Name, p.Dir, ideName})
	}
	Success(fmt.Sprintf("找到 %d 个项目：", len(projects)))
	fmt.Print(renderTable([]string{"ID", "项目名称", "项目目录", "默认 IDE"}, rows))
}

// PrintIDEs 打印 IDE 列表。
func PrintIDEs() {
	ides := store.FindAllIDEs()
	if len(ides) == 0 {
		Info("没有找到任何 IDE。请使用 'ordo ide add' 注册一个 IDE。")
		return
	}
	rows := make([][]string, 0, len(ides))
	for _, ide := range ides {
		rows = append(rows, []string{ide.ID, ide.Name, ide.Path})
	}
	Success(fmt.Sprintf("找到 %d 个 IDE：", len(ides)))
	fmt.Print(renderTable([]string{"ID", "IDE 名称", "安装路径"}, rows))
}
