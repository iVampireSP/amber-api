package entity

import (
	"time"
)

type File struct {
	Model

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
