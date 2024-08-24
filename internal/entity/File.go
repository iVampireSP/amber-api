package entity

import (
	"time"
)

// File 的 id 为 string 类型，不是 int64。主键为 AUTO RANDOM
// 事实上之后我们也要改变其他 entity 的 Id 为 string，所以从这里开始实验更改。
type File struct {
	Entity    `xorm:"extends"`
	Url       *string    `xorm:"varchar(255) notnull" json:"url"`
	UrlHash   *string    `xorm:"varchar(255) notnull" json:"url_hash"`
	FileHash  string     `xorm:"varchar(255) notnull" json:"file_hash"`
	MimeType  string     `xorm:"mime_type int(11) notnull" json:"mime_type"`
	Path      string     `xorm:"varchar(255) notnull" json:"path"`
	ExpiredAt *time.Time `xorm:"TIMESTAMP null" json:"expired_at"`
}

func (a *File) TableName() string {
	return "files"
}
