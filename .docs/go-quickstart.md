# Go 快速上手指南

## 前提

Go 1.26 已安装，本机已有编译好的二进制。

## 每天怎么用

项目在 `/Users/owocc/projects/opensource/ordo`，直接在终端：

```sh
cd ~/projects/opensource/ordo

# 看已管理的项目
./dist/bin/ordo ls

# 看已注册的 IDE
./dist/bin/ordo ide ls

# 添加项目（缺参数自动交互）
./dist/bin/ordo add myapp ~/projects/myapp

# 打开项目
./dist/bin/ordo open myapp

# 启动 TUI 界面
./dist/bin/ordo tui

# 查看从 agent 发现的项目（opencode / codex / claude）
./dist/bin/ordo agent ls
```

### 安装到全局（之后可直接 `ordo`）

```sh
make install
# 或手动:
sudo cp dist/bin/ordo /usr/local/bin/
```

## 如果你要改代码

编译运行：

```sh
make build     # 编译到 dist/bin/ordo
make run       # 编译 + 运行一次
make dev       # go run 直接跑（不保存二进制，适合调试）
make tui       # 编译 + 启动 TUI
```

代码改了重新编译即可，Go 编译很快（几秒）。

## 如果你要重新全平台编译

```sh
make release   # 同时编译 darwin/linux/windows × amd64/arm64
```

产物在 `dist/` 目录下。

## 记住

所有 `go` 命令需要代理（已写在 Makefile 里自动应用），你直接敲 `make xxxx` 就行，不需要手动设。

如果自己敲 `go build`，需要：

```sh
GOPROXY=https://goproxy.cn,direct GOSUMDB=off go build ./cmd/ordo
```
