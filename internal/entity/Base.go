package entity

type Base struct {
	ID        int64  `xorm:"id BigInt notnull pk" json:"id,string"`
	CreatedAt string `xorm:"created timestamp" json:"created_at"`
	UpdatedAt string `xorm:"updated timestamp" json:"updated_at"`
}
