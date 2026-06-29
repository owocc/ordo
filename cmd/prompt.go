package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/owocc/ordo/internal/model"
	"github.com/owocc/ordo/internal/store"
)

var reader = bufio.NewReader(os.Stdin)

// promptText 读取一行用户输入，去掉首尾空白。
func promptText(prompt string) string {
	fmt.Print(prompt + " ")
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// promptConfirm 询问 y/n，默认 n。
func promptConfirm(prompt string) bool {
	s := promptText(prompt + " [y/N]")
	return strings.EqualFold(s, "y") || strings.EqualFold(s, "yes")
}

// promptSelectProjects 列出所有项目让用户用逗号分隔多选（或输入 all）。
func promptSelectProjects(projects []model.Project) []model.Project {
	fmt.Println("请选择要打开的项目（逗号分隔序号，或输入 all）：")
	for i, p := range projects {
		fmt.Printf("  %d. %s  (%s)\n", i+1, p.Name, p.Dir)
	}
	for {
		s := strings.ToLower(promptText("输入序号："))
		if s == "" {
			return nil
		}
		if s == "all" {
			return append([]model.Project{}, projects...)
		}
		var out []model.Project
		ok := true
		for _, part := range strings.Split(s, ",") {
			part = strings.TrimSpace(part)
			n, err := strconv.Atoi(part)
			if err != nil || n < 1 || n > len(projects) {
				ok = false
				break
			}
			out = append(out, projects[n-1])
		}
		if !ok {
			fmt.Println("无效输入，请重新选择。")
			continue
		}
		return out
	}
}

// promptSelectIDE 列出已注册的 IDE 让用户选择；返回选中的 IDE id。
// allowNone=true 时提供"无 IDE"选项；selectionEmpty 时返回""。
// 默认 defaultID 为预选项（直接回车选用）。
func promptSelectIDE(prompt string, defaultID string) string {
	ides := store.FindAllIDEs()
	if len(ides) == 0 {
		return ""
	}
	fmt.Println(prompt)
	for i, ide := range ides {
		marker := "  "
		if ide.ID == defaultID {
			marker = "> "
		}
		desc := ide.Desc
		if desc != "" {
			desc = "  (" + desc + ")"
		}
		fmt.Printf("%s%d. %s%s\n", marker, i+1, ide.Name, desc)
	}
	fmt.Printf("  0. 无 IDE\n")
	for {
		s := promptText("输入序号：")
		if s == "" && defaultID != "" {
			return defaultID
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 || n > len(ides) {
			fmt.Println("无效输入，请重新选择。")
			continue
		}
		if n == 0 {
			return ""
		}
		return ides[n-1].ID
	}
}
