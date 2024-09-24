package entity

import (
	"strings"
	"time"
)

type File struct {
	Model

	Url      *string `json:"url"`
	UrlHash  *string `json:"url_hash"`
	FileHash string  `json:"file_hash"`
	MimeType string  `json:"mime_type"`
	Path     string  `json:"-"`
	Size     int64   `json:"size"`
	//Public   bool    `json:"public"` // 是否公开，访客上传的文件应始终公开，或归属于所有者
	// TODO: 移除 file 的到期时间，如果当 file 没有任何引用的时候再删除
	// 因为有外键，所以直接删除是删不掉的，必须删除消息
	ExpiredAt *time.Time `json:"expired_at"`
}

func (f *File) TableName() string {
	return "files"
}

func (f *File) IsImage() bool {
	// 如果 mimetype 不是 image/  开头
	return strings.HasPrefix(f.MimeType, "image/")
}
