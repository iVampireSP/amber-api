package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "4",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE chats (
  id   bigint unsigned AUTO_INCREMENT,
  name varchar(255) DEFAULT NULL,
  assistant_id bigint unsigned NOT NULL,
  user_id bigint unsigned DEFAULT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`

			_, err := tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			rawSQL = `CREATE INDEX chats_assistant_id_index ON chats (assistant_id);`
			_, err = tx.Exec(rawSQL)
			if err != nil {
				return err
			}
			rawSQL = `CREATE INDEX chats_user_id_index ON chats (user_id);`
			_, err = tx.Exec(rawSQL)
			if err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// drop table
			_, err := tx.Exec("DROP TABLE IF EXISTS chats;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
