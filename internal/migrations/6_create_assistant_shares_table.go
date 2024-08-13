package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "6",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE assistant_shares (
  id   bigint unsigned AUTO_INCREMENT,
  assistant_id bigint unsigned NOT NULL,
  token varchar(255) NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

`

			_, err := tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`CREATE INDEX assistant_shares_assistant_id_index ON assistant_shares (assistant_id);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`CREATE INDEX assistant_shares_token_index ON assistant_shares (token);`)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// drop table
			_, err := tx.Exec("DROP TABLE IF EXISTS assistant_shares;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
