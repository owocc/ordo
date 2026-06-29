package agent

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type opencode struct{}

func (*opencode) Name() string { return "opencode" }

func (*opencode) Discover() ([]Project, error) {
	dbPath := opencodeDBPath()
	if dbPath == "" {
		return nil, nil
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, nil
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT directory, COUNT(*) as cnt, MAX(time_updated) as last
		FROM session
		WHERE directory != ''
		GROUP BY directory
		ORDER BY last DESC
	`)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var out []Project
	for rows.Next() {
		var dir string
		var cnt int
		var ts int64
		if err := rows.Scan(&dir, &cnt, &ts); err != nil {
			continue
		}
		if !isDir(dir) {
			continue
		}
		out = append(out, Project{
			Path:         dir,
			Sources:      []string{"opencode"},
			SessionCount: cnt,
			LastActive:   time.UnixMilli(ts),
		})
	}
	return out, nil
}

func opencodeDBPath() string {
	// opencode 数据目录: 常见位置 ~/.local/share/opencode/
	if home, err := os.UserHomeDir(); err == nil {
		paths := []string{
			filepath.Join(home, ".local", "share", "opencode", "opencode.db"),
		}
		// 也尝试 XDG_DATA_HOME
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg != "" {
			paths = append([]string{filepath.Join(xdg, "opencode", "opencode.db")}, paths...)
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}
