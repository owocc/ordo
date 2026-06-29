package model

import (
	"fmt"
	"regexp"
	"time"
)

var nameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateName 校验标识符名称：只能包含字母、数字、连字符和下划线。
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("名称不能为空")
	}
	if !nameRe.MatchString(name) {
		return fmt.Errorf("名称只能包含字母、数字、连字符(-)和下划线(_)，不能包含空格")
	}
	return nil
}

// DisplayOrName 返回显示名称，没有设则用 Name。
func (p Project) DisplayOrName() string {
	if p.DisplayName != "" {
		return p.DisplayName
	}
	return p.Name
}

// DisplayOrName 返回显示名称，没有设则用 Name。
func (ide IDE) DisplayOrName() string {
	if ide.DisplayName != "" {
		return ide.DisplayName
	}
	return ide.Name
}

// Project 表示一个被 ordo 管理的项目。
type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	DisplayName   string    `json:"displayName,omitempty"`
	Dir           string    `json:"dir"`
	Tags          []string  `json:"tag,omitempty"`
	RelationIDEID string    `json:"relationIdeId,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
}

// IDE 表示一个已注册的 IDE。
type IDE struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Path        string `json:"path"`
	Args        string `json:"args,omitempty"`
	Desc        string `json:"desc,omitempty"`
}

// Config 是持久化到磁盘的整体配置结构（与旧 Deno 版兼容）。
type Config struct {
	Projects   []Project `json:"projects"`
	IDEs       []IDE     `json:"ides"`
	DefaultIDE *IDE      `json:"defaultIde,omitempty"`
}
