package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "1",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE tools (
  id   bigint unsigned AUTO_INCREMENT,
  name varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  description varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  discovery_url varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  api_key varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  data json DEFAULT NULL,
  user_id bigint NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
`
			_, err := tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			rawSQL = `CREATE INDEX tools_user_id_index ON tools (user_id);`
			_, err = tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// drop table
			_, err := tx.Exec("DROP TABLE IF EXISTS tools;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
