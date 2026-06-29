# Agent 项目发现功能方案

## 目标

从本机 agent 应用（opencode、codex、claude code）的存储目录中扫描出项目目录，方便用不同编辑器打开。

## 三个来源的存储布局

| 来源       | 路径                                  | 数据格式          | 如何获取项目路径                                                                 |
| ---------- | ------------------------------------- | ----------------- | -------------------------------------------------------------------------------- |
| Claude Code | `~/.claude/projects/<编码路径>/`      | 子目录名即 cwd    | 目录名 `-Users-owocc-projects-foo` → 去掉前缀 `-`，剩余 `-` 替换为 `/` → `/Users/owocc/projects/foo` |
| Codex      | `~/.codex/sessions/YYYY/MM/DD/rollout-*.jsonl` | JSONL 首行        | 只读首行，`session_meta.payload.cwd` 字段                                          |
| OpenCode   | `~/.local/share/opencode/opencode.db` | SQLite            | `SELECT directory FROM session GROUP BY directory`                              |

## 数据结构

```go
type Project struct {
    Path         string      // 项目目录绝对路径
    Sources      []string    // 来源标识（claude / codex / opencode）
    SessionCount int         // 该目录下的会话总数
    LastActive   time.Time   // 最近一次会话时间（用于排序）
}
```

## 聚合逻辑

- 所有 Source 并发执行 `Discover()`，结果按 `Path` 去重合并
- Sources 取并集、SessionCount 累加、LastActive 取最大值
- 最终按 `LastActive` 降序排列
- 过滤掉 `Path == "/"` 或空路径

## 命令设计

- `ordo agent ls` — 列出所有发现的项目（带来源徽标、会话数、最近活跃时间）
- `ordo agent open <#序号|路径>` — 交互式选 IDE 后打开项目
- `--source` 过滤（预留，当前未实现）

## TUI 集成

新增 `viewAgentProjects` 视图，Tab 在 项目列表 → IDE 列表 → Agent 列表 间循环。
Agent 列表中按 Enter 用第一个已注册 IDE 打开项目。

## 依赖

- `modernc.org/sqlite`（纯 Go SQLite，无 CGO，支持交叉编译）
