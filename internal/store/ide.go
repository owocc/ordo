package store

import (
	"github.com/google/uuid"
	"github.com/owocc/ordo/internal/model"
)

// FindAllIDEs 返回所有已注册的 IDE。
func FindAllIDEs() []model.IDE {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	out := make([]model.IDE, len(cfg.IDEs))
	copy(out, cfg.IDEs)
	return out
}

// FindOneIDE 按 id 或 name 查找 IDE。
func FindOneIDE(idOrName string) *model.IDE {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	for i := range cfg.IDEs {
		ide := &cfg.IDEs[i]
		if ide.ID == idOrName || ide.Name == idOrName {
			cp := *ide
			return &cp
		}
	}
	return nil
}

// AddIDE 注册一个新 IDE，id 由 store 生成。
func AddIDE(ide model.IDE) (model.IDE, error) {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	ide.ID = uuid.NewString()
	cfg.IDEs = append(cfg.IDEs, ide)
	if err := persist(); err != nil {
		cfg.IDEs = cfg.IDEs[:len(cfg.IDEs)-1]
		return model.IDE{}, err
	}
	return ide, nil
}

// DeleteIDE 按 id 或 name 删除 IDE。
func DeleteIDE(idOrName string) (int, error) {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	kept := cfg.IDEs[:0]
	removed := 0
	for _, ide := range cfg.IDEs {
		if ide.ID == idOrName || ide.Name == idOrName {
			removed++
			continue
		}
		kept = append(kept, ide)
	}
	cfg.IDEs = kept
	return removed, persist()
}
