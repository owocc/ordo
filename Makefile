# 修改这些变量以控制构建
BINARY    := ordo
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
DIST      := dist
LDFLAGS   := -s -w
GOFLAGS   := -trimpath
TARGETS   := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# 国内代理（按需修改或注释掉）
export GOPROXY ?= https://goproxy.cn,direct
export GOSUMDB ?= off

.PHONY: build run dev test fmt vet clean release install uninstall tui help

build: ## 编译当前平台的二进制到 dist/bin/ordo
	@mkdir -p $(DIST)/bin
	go build $(GOFLAGS) -ldflags '$(LDFLAGS) -X main.version=$(VERSION)' -o $(DIST)/bin/$(BINARY) ./cmd/ordo

run: build ## 运行一次 CLI
	$(DIST)/bin/$(BINARY)

dev: ## 运行 CLI（go run，便于调试）
	go run ./cmd/ordo

tui: build ## 进入交互式 TUI
	$(DIST)/bin/$(BINARY) tui

release: ## 交叉编译多平台二进制到 dist/
	@mkdir -p $(DIST)
	@for target in $(TARGETS); do \
		os=$${target%/*}; arch=$${target#*/}; \
		ext=""; if [ $$os = windows ]; then ext=.exe; fi; \
		echo "→  $$os/$$arch"; \
		GOOS=$$os GOARCH=$$arch go build $(GOFLAGS) -ldflags '$(LDFLAGS)' \
			-o $(DIST)/$(BINARY)-$$os-$$arch$$ext ./cmd/ordo; \
	done
	@echo "✓ 构建完成: $(DIST)/"

install: build ## 安装到 $$GOPATH/bin
	@cp $(DIST)/bin/$(BINARY) "$$(go env GOPATH)/bin/"

uninstall: ## 卸载
	@rm -f "$$(go env GOPATH)/bin/$(BINARY)"

fmt: ## go fmt
	@go fmt ./...

vet: ## go vet
	@go vet ./...

test: vet ## go test（项目当前无测试用例）
	@go test ./...

clean: ## 清理 dist
	@rm -rf $(DIST)

help: ## 显示本帮助
	@awk 'BEGIN {FS = ":.*##"; printf "Make 命令：\n"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)