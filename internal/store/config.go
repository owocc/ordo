package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/owocc/ordo/internal/model"
)

// 兼容旧 Deno 版 npm:conf 的存储位置：
//
//	macOS   ~/Library/Preferences/ordo/config.json
//	Linux   ~/.config/ordo/config.json
//	Windows %APPDATA%\ordo\config.json
func configPath() string {
	var base string
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, "Library", "Preferences", "ordo")
	case "windows":
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			home, _ := os.UserHomeDir()
			appdata = filepath.Join(home, "AppData", "Roaming")
		}
		base = filepath.Join(appdata, "ordo")
	default:
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			home, _ := os.UserHomeDir()
			xdg = filepath.Join(home, ".config")
		}
		base = filepath.Join(xdg, "ordo")
	}
	return filepath.Join(base, "config.json")
}

var (
	mu   sync.Mutex
	once sync.Once
	cfg  *model.Config
	path string
)

func initStore() {
	once.Do(func() {
		path = configPath()
		cfg = load()
	})
}

func load() *model.Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return &model.Config{
			Projects: []model.Project{},
			IDEs:     []model.IDE{},
		}
	}
	var c model.Config
	if err := json.Unmarshal(data, &c); err != nil {
		return &model.Config{}
	}
	if c.Projects == nil {
		c.Projects = []model.Project{}
	}
	if c.IDEs == nil {
		c.IDEs = []model.IDE{}
	}
	return &c
}

func persist() error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Raw 返回当前内存中的配置副本，供 ui/tui 使用。
func Raw() *model.Config {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	cp := *cfg
	return &cp
}
