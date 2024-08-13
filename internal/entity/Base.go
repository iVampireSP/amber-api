package entity

type Base struct {
	Id        int64  `json:"id"`
	CreatedAt string `xorm:"created timestamp" json:"created_at"`
	UpdatedAt string `xorm:"updated timestamp" json:"updated_at"`
}
