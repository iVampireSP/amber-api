package migrations

import (
	"fmt"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "1",
		Migrate: func(tx *xorm.Engine) error {
			fmt.Println("up")
			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			fmt.Println("down")
			return nil
		},
	})
}
