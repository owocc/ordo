package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/owocc/ordo/internal/model"
)

// FindAllProjects 返回所有项目。
func FindAllProjects() []model.Project {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	out := make([]model.Project, len(cfg.Projects))
	copy(out, cfg.Projects)
	return out
}

// FindOneProject 按 id 或 name 查找项目。
func FindOneProject(idOrName string) *model.Project {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	for i := range cfg.Projects {
		p := &cfg.Projects[i]
		if p.ID == idOrName || p.Name == idOrName {
			cp := *p
			return &cp
		}
	}
	return nil
}

// FindProjectAndIDE 返回匹配的项目及其关联的 IDE。
func FindProjectAndIDE(idOrNames []string) []model.Project {
	// 通过 GetProjectsAndIDE 的非写入版返回结构，但这里为了兼容
	// 直接返回 project 副本，IDE 查询交给上层。
	matches := []model.Project{}
	for _, name := range idOrNames {
		if p := FindOneProject(name); p != nil {
			matches = append(matches, *p)
		}
	}
	return matches
}

// AddProject 添加新项目，id 由 store 生成。
func AddProject(p model.Project) (model.Project, error) {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	p.ID = uuid.NewString()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	cfg.Projects = append(cfg.Projects, p)
	if err := persist(); err != nil {
		cfg.Projects = cfg.Projects[:len(cfg.Projects)-1]
		return model.Project{}, err
	}
	return p, nil
}

// DeleteProject 按 id 或 name 删除项目。
func DeleteProject(idOrName string) (int, error) {
	initStore()
	mu.Lock()
	defer mu.Unlock()
	kept := cfg.Projects[:0]
	removed := 0
	for _, p := range cfg.Projects {
		if p.ID == idOrName || p.Name == idOrName {
			removed++
			continue
		}
		kept = append(kept, p)
	}
	cfg.Projects = kept
	// 若用户传入名字则原删除逻辑允许找不到也算成功（返回 0）
	return removed, persist()
}
