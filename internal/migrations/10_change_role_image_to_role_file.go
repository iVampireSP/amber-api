package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "10",
		Migrate: func(tx *xorm.Engine) error {
			var err error
			var sql = `
UPDATE chat_messages SET role = 'file' WHERE role = 'image';
`
			_, err = tx.Exec(sql)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			_, err := tx.Exec(`UPDATE chat_messages SET role = 'image' WHERE role = 'file';`)
			if err != nil {
				return err
			}

			return nil
		},
	})
}
