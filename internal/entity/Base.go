package entity

type Base struct {
	ID        int    `xorm:"id int(11) notnull pk autoincr" json:"id"`
	CreatedAt string `xorm:"created timestamp" json:"created_at"`
	UpdatedAt string `xorm:"updated timestamp" json:"updated_at"`
}
