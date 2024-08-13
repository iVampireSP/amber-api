package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "3",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE assistant_tools (
 id  bigint AUTO_RANDOM,
 assistant_id bigint NOT NULL,
 tool_id bigint NOT NULL,
 created_at timestamp NULL DEFAULT NULL,
 updated_at timestamp NULL DEFAULT NULL,
 PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`

			_, err := tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			rawSQL = `CREATE INDEX assistant_tools_tool_id_index ON assistant_tools (tool_id);`
			_, err = tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			rawSQL = `CREATE INDEX assistant_tools_assistant_id_index ON assistant_tools (assistant_id);`
			_, err = tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// drop table
			_, err := tx.Exec("DROP TABLE IF EXISTS assistant_tools;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
