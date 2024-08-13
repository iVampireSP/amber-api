package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "5",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE chat_messages (
  id  bigint AUTO_RANDOM,
  chat_id bigint NOT NULL,
  content text NOT NULL,
  role varchar(255) DEFAULT NULL,
  prompt_tokens bigint NOT NULL DEFAULT 0,
  completion_tokens bigint NOT NULL DEFAULT 0,
  total_tokens bigint NOT NULL DEFAULT 0,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- drop primary key
`

			_, err := tx.Exec(rawSQL)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`CREATE INDEX chat_messages_chat_id_index ON chat_messages (chat_id);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`CREATE INDEX chat_messages_created_at_index ON chat_messages (created_at);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`CREATE INDEX chat_messages_role_index ON chat_messages (role);`)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			// drop table
			_, err := tx.Exec("DROP TABLE IF EXISTS chat_messages;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
