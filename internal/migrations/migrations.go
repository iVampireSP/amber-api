package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

var migrations []*migrate.Migration
var m *migrate.Migrate

func NewMigrate(x *xorm.Engine) {
	m = migrate.New(x, &migrate.Options{
		TableName:    "migrations",
		IDColumnName: "id",
	}, migrations)
}

func Migrate() {
	err := m.Migrate()
	if err != nil {
		panic(err)
		return
	}
}

func Rollback() {
	err := m.RollbackLast()
	if err != nil {
		panic(err)
		return
	}
}
