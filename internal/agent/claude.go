package agent

import (
	"os"
	"path/filepath"
	"strings"
)

type claude struct{}

func (*claude) Name() string { return "claude" }

func (*claude) Discover() ([]Project, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".claude", "projects")
	entries, err := os.ReadDir(dir)
	if err != nil {
		// 未安装 claude
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Project
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := decodeClaudeDir(e.Name())
		if path == "" || !isDir(path) {
			continue
		}
		sessions := countJSONLs(filepath.Join(dir, e.Name()))
		if sessions == 0 {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, Project{
			Path:         path,
			Sources:      []string{"claude"},
			SessionCount: sessions,
			LastActive:   fi.ModTime(),
		})
	}
	return out, nil
}

// decodeClaudeDir 将 "-Users-owocc-projects-foo" 还原为 "/Users/owocc/projects/foo"。
// Claude 编码: 去掉前导 "-"，剩余 "-" 替换为 "/"。
func decodeClaudeDir(name string) string {
	if !strings.HasPrefix(name, "-") {
		return ""
	}
	return "/" + strings.ReplaceAll(name[1:], "-", "/")
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

func countJSONLs(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			n++
		}
	}
	return n
}
