package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "8",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
-- add disable_default_prompt(bool) column
ALTER TABLE "public"."assistants" ADD COLUMN "disable_default_prompt" bool NOT NULL DEFAULT false;
`

			_, err := tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// drop table
			_, err := tx.Exec(`
-- Drop column
ALTER TABLE "public"."assistants" DROP COLUMN "disable_default_prompt";
`)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
