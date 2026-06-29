package model

import "time"

// Project 表示一个被 ordo 管理的项目。
type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Dir           string    `json:"dir"`
	Tags          []string  `json:"tag,omitempty"`
	RelationIDEID string    `json:"relationIdeId,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
}

// IDE 表示一个已注册的 IDE。
type IDE struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Args string `json:"args,omitempty"`
	Desc string `json:"desc,omitempty"`
}

// Config 是持久化到磁盘的整体配置结构（与旧 Deno 版兼容）。
type Config struct {
	Projects   []Project `json:"projects"`
	IDEs       []IDE     `json:"ides"`
	DefaultIDE *IDE      `json:"defaultIde,omitempty"`
}
