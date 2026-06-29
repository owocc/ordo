package opener

import (
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/owocc/ordo/internal/model"
)

// 预测 IDE 启动命令：
//   - 旧 Deno 版逻辑：darwin/linux 在绝对路径前加 "open"；相对路径/命令原样使用。
//   - Windows 不加前缀。
func getPrefix() string {
	// 注：原版对 linux 也用 "open"，但 linux 通常无 open 命令。这里维持原版行为
	// 以兼容旧配置；用户通常用 code/idea 等相对命令而非绝对路径，不会触发此前缀。
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return "open"
	}
	return ""
}

// BuildCommand 为单个项目构造启动命令字符串。
func BuildCommand(p model.Project, ide *model.IDE) string {
	idePath := ide.Path
	if idePath != "" && filepath.IsAbs(idePath) {
		// 用引号包裹绝对路径以处理空格
		idePath = getPrefix() + " \"" + idePath + "\""
	}
	dir := p.Dir
	if dir == "" {
		return idePath
	}
	return idePath + " " + dir
}

// Run 直接通过系统 shell 执行命令（兼容旧版 exec 字符串行为）。
func Run(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", command)
	}
	return exec.Command("sh", "-c", command)
}
