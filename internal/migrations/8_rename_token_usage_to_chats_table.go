package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "8",
		Migrate: func(tx *xorm.Engine) error {
			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			return nil
		},
	})
}
