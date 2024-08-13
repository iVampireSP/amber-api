package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "2",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE assistants (
  id serial NOT NULL primary key ,
  name varchar(255) DEFAULT NULL,
  description varchar(255) DEFAULT NULL,
  prompt text DEFAULT NULL,
  disable_default_prompt boolean NOT NULL,
  user_id bigint NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`

			_, err := tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			rawSQL = `CREATE INDEX assistants_user_id_index ON assistants (user_id);`
			_, err = tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// drop table
			_, err := tx.Exec("DROP TABLE IF EXISTS assistants;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
