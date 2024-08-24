package migrations

import (
	"rag-new/internal/schema"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

type ChatMessage struct {
	Id      int64
	Content string          `xorm:"varchar(255) notnull"`
	Role    schema.ChatRole `xorm:"varchar(255) notnull"`
	Hidden  bool            `xorm:"bool notnull"`
}

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "9",
		Migrate: func(tx *xorm.Engine) error {
			var err error
			var sql = `
ALTER TABLE chat_messages ADD COLUMN hidden boolean NOT NULL DEFAULT false AFTER role;
`
			_, err = tx.Exec(sql)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index chat_messages_hidden_index on chat_messages (hidden);`)
			if err != nil {
				return err
			}
			sql = `update chat_messages set hidden = true where role LIKE "%_hide"`
			_, err = tx.Query(sql)
			if err != nil {
				return err
			}

			//sql = `update chat_messages set role = "system" where role = "system_hide"`
			//_, err = tx.Query(sql)
			//if err != nil {
			//	return err
			//}
			//
			//sql = `update chat_messages set role = "user" where role = "user_hide"`
			//_, err = tx.Query(sql)
			//if err != nil {
			//	return err
			//}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			_, err := tx.Exec(`drop index chat_messages_hidden_index on chat_messages;`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`ALTER TABLE chat_messages DROP COLUMN hidden;`)
			if err != nil {
				return err
			}

			return nil
		},
	})
}
