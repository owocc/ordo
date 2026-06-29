package agent

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type codex struct{}

func (*codex) Name() string { return "codex" }

func (*codex) Discover() ([]Project, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	base := filepath.Join(home, ".codex", "sessions")
	if _, err := os.Stat(base); os.IsNotExist(err) {
		return nil, nil
	}
	type hit struct {
		cwd  string
		mtime time.Time
	}
	var hits []hit

	// 扫描 rollout-*.jsonl 文件
	filepath.Walk(base, func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() || !strings.HasPrefix(fi.Name(), "rollout-") || !strings.HasSuffix(p, ".jsonl") {
			return nil
		}
		cwd, mtime := readCodexCWD(p, fi.ModTime())
		if cwd != "" {
			hits = append(hits, hit{cwd, mtime})
		}
		return nil
	})
	if len(hits) == 0 {
		return nil, nil
	}

	// 按 cwd 聚合
	type agg struct {
		count int
		last  time.Time
	}
	aggMap := map[string]*agg{}
	for _, h := range hits {
		a, ok := aggMap[h.cwd]
		if !ok {
			a = &agg{}
			aggMap[h.cwd] = a
		}
		a.count++
		if h.mtime.After(a.last) {
			a.last = h.mtime
		}
	}
	out := make([]Project, 0, len(aggMap))
	for cwd, a := range aggMap {
		out = append(out, Project{
			Path:         cwd,
			Sources:      []string{"codex"},
			SessionCount: a.count,
			LastActive:   a.last,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastActive.After(out[j].LastActive)
	})
	return out, nil
}

type codexMeta struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type codexMetaPayload struct {
	Cwd string `json:"cwd"`
}

func readCodexCWD(path string, fallback time.Time) (string, time.Time) {
	f, err := os.Open(path)
	if err != nil {
		return "", fallback
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	// 首行可能超过 64KB（base_instructions 很长），但我们要的 cwd 在前部
	// 设置大 buffer 以防截断
	sc.Buffer(make([]byte, 512*1024), 512*1024)
	if !sc.Scan() {
		return "", fallback
	}
	line := sc.Bytes()
	var meta codexMeta
	if err := json.Unmarshal(line, &meta); err != nil {
		return "", fallback
	}
	if meta.Type != "session_meta" {
		return "", fallback
	}
	var pl codexMetaPayload
	if err := json.Unmarshal(meta.Payload, &pl); err != nil {
		return "", fallback
	}
	if pl.Cwd == "" || !isDir(pl.Cwd) {
		return "", fallback
	}
	return pl.Cwd, fallback
}
