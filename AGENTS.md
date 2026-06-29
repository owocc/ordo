# AGENTS.md

Ordo 现已是 **Go** 项目（旧 Deno/TS 代码保留在原位但已不用）。模块名 `github.com/owocc/ordo`，最低 Go 1.26。

## 构建 / 运行

⚠️ 网络依赖：`go mod` 拉取需要代理，所有 `go` 命令必须带：

```sh
export GOPROXY=https://goproxy.cn,direct GOSUMDB=off
```

（无此设置 `go mod tidy`/`go build` 会因 sumdb 超时失败。`Makefile` 默认已 export。）

常用命令（见 `Makefile`）：

- `make build` — 编译当前平台二进制到 `dist/bin/ordo`
- `make run` / `make dev` — 一次运行 / `go run`
- `make tui` — 进入交互式 TUI 模式
- `make release` — 交叉编译 5 个目标（darwin amd64/arm64、linux amd64/arm64、windows amd64）到 `dist/`
- `make install` — 装到 `$GOPATH/bin`
- `make fmt` / `make vet` / `make test` — 无测试用例，`test` 等同 `vet`

未配置 CI、lint。无测试。

## 架构

- **入口** `cmd/ordo/main.go` → `cmd` 包（cobra 命令：root/add/open/ls/rm/ide/tui）。
- **模型** `internal/model`：`Project` / `IDE` / `Config`，JSON tag 与旧 Deno 版兼容。
- **存储** `internal/store`：JSON 持久化在 OS 配置目录，路径由 `config.go` 按 OS 分发（macOS `~/Library/Preferences/ordo/config.json`，与旧 `npm:conf` 兼容；Linux `~/.config/ordo/`；Windows `%APPDATA%\ordo\`）。`store` 包内带 `sync.Mutex`，所有读写经此，**不要** 直接读写文件。`store.FindOneProject` / `FindOneIDE` 接受 `id` 或 `name`。
- **打开项目** `internal/opener`：构造命令字符串，维持旧逻辑——仅在 IDE 路径是绝对路径时给 darwin/linux 加 `open` 前缀（Linux 上 `open` 并非标准命令，故使用相对命令/`code`/`idea` 时不触发该前缀）。`opener.Run` 返回 `*exec.Cmd`，经 `sh -c`（windows `cmd /c`）执行——因为命令是拼接字符串而非 argv。
- **CLI 表格** `internal/ui`：无边框、列间三空格，用 lipgloss 着色，复刻旧 `cli-table3` 的 `blankBorder` 风格。
- **TUI** `internal/tui`：bubbletea `Model` 按 `view` 状态机切视图（项目列表 / IDE 列表 / 添加表单 / 删除确认）。Enter 打开、d 删除、a 添加项目、A 添加 IDE、Tab 换页、q 退出。表单用 `bubbles/textinput`。
- **旧 Deno 实现**：`cli/`、`store/*.store.ts`、`lib/`、`types/`、`utils/`、`web-ui/`、`deno.json`、`deno.lock` 仍在仓库但不再被使用；改任何功能请只动 Go 代码。

## 约定

- Go 风格：标准 `gofmt`，4 空格无关（tab）；中文注释与文案 OK。
- 用户面消息走 `internal/ui` 的 `Success/Info/Warn/Error`（彩色 ✓/!）。CLI 命令失败用 `cmd.fail()` 返回 error 让 cobra 退出码非 0。
- 交互式输入（缺参时）走 `cmd/prompt.go` 的 `promptText`/`promptSelectIDE`（从 stdin 读行）；TUI 内的输入走 `textinput`。
- 新增依赖：先 `go get` 再在源码 import；保持 `go.mod` 与 `go.sum` 同步（用 `make` 触发的 `go mod` 会自动用代理）。
- 不要触碰旧 `.ts`/`.tsx` 文件除非要迁移其逻辑到 Go。