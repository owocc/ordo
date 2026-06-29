package agent

import (
	"sort"
	"time"

	"github.com/owocc/ordo/internal/model"
)

// Project 是从 agent store 中发现的项目目录。
type Project struct {
	Path         string    // 项目目录绝对路径
	Sources      []string  // 来源: claude / codex / opencode
	SessionCount int       // 该目录的总会话数
	LastActive   time.Time // 最近会话时间
}

// Source 是 agent 存储源（claude / codex / opencode）的统一接口。
type Source interface {
	Name() string
	Discover() ([]Project, error)
}

// defaultSources 返回所有内置的 source。
func defaultSources() []Source {
	return []Source{&claude{}, &codex{}, &opencode{}}
}

// Discover 聚合所有 source 的结果，按 path 去重合并。
func Discover() ([]Project, error) {
	return DiscoverSources(defaultSources())
}

// DiscoverSources 接受自定义 sources 列表，用于测试或 --source 过滤。
func DiscoverSources(sources []Source) ([]Project, error) {
	type result struct {
		projects []Project
		err      error
		srcName  string
	}
	ch := make(chan result, len(sources))
	for _, s := range sources {
		go func(s Source) {
			projects, err := s.Discover()
			ch <- result{projects, err, s.Name()}
		}(s)
	}

	byPath := map[string]*Project{}
	for range sources {
		r := <-ch
		// 单源失败不阻断整体
		_ = r.err
		for _, p := range r.projects {
			if p.Path == "/" || p.Path == "" {
				continue
			}
			if existing, ok := byPath[p.Path]; ok {
				existing.SessionCount += p.SessionCount
				existing.Sources = append(existing.Sources, p.Sources...)
				if p.LastActive.After(existing.LastActive) {
					existing.LastActive = p.LastActive
				}
			} else {
				cp := p
				byPath[p.Path] = &cp
			}
		}
	}
	out := make([]Project, 0, len(byPath))
	for _, p := range byPath {
		out = append(out, *p)
	}
	// 最近活跃倒序
	sortByLastActiveDesc(out)
	return out, nil
}

// ToProject 将 agent.Project 转为 opener 可用的 model.Project。
func ToProject(p *Project) model.Project {
	return model.Project{
		Name: p.Path,
		Dir:  p.Path,
	}
}

// sortByLastActiveDesc 按 LastActive 降序排序。
func sortByLastActiveDesc(projects []Project) {
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastActive.After(projects[j].LastActive)
	})
}
