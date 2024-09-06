package entity

import (
	"time"
)

type File struct {
	Model

	Url      *string `json:"url"`
	UrlHash  *string `json:"url_hash"`
	FileHash string  `json:"file_hash"`
	MimeType string  `json:"mime_type"`
	Path     string  `json:"path"`
	Size     int64   `json:"size"`
	// 可见性，如果是 true 代表公开，可以直接下载。否则则是私有，需要临时密钥
	Public bool `json:"public"`
	// TODO: 移除 file 的到期时间，如果当 file 没有任何引用的时候再删除
	// 因为有外键，所以直接删除是删不掉的，必须删除消息
	ExpiredAt *time.Time `json:"expired_at"`
}

func (a *File) TableName() string {
	return "files"
}
