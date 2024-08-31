package entity

import (
	"time"
)

type File struct {
	Model

	Url      *string `xorm:"varchar(255) notnull" json:"url"`
	UrlHash  *string `xorm:"varchar(255) notnull" json:"url_hash"`
	FileHash string  `xorm:"varchar(255) notnull" json:"file_hash"`
	MimeType string  `xorm:"mime_type int(11) notnull" json:"mime_type"`
	Path     string  `xorm:"varchar(255) notnull" json:"path"`
	// TODO: 移除 file 的到期时间，如果当 file 没有任何引用的时候再删除
	// 因为有外键，所以直接删除是删不掉的，必须删除消息
	ExpiredAt *time.Time `xorm:"TIMESTAMP null" json:"expired_at"`
}

func (a *File) TableName() string {
	return "files"
}
