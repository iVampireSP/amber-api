package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "7",
		Migrate: func(tx *xorm.Engine) error {
			_, err := tx.Exec(`ALTER TABLE chats add column expired_at timestamp NULL DEFAULT NULL AFTER user_id;`)
			if err != nil {
				return err
			}
			_, err = tx.Exec(`ALTER TABLE chats add column owner varchar(255) DEFAULT NULL AFTER user_id;`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`ALTER TABLE chats add column guest_id varchar(255) DEFAULT NULL AFTER user_id;`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index chats_expired_at_index on chats (expired_at);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index chats_owner_index on chats (owner);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index chats_guest_id_index on chats (guest_id);`)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			_, err := tx.Exec("DROP INDEX chats_owner_index on chats;")
			if err != nil {
				return err
			}

			_, err = tx.Exec(`DROP INDEX chats_expired_at_index on chats;`)
			if err != nil {
				return err
			}
			_, err = tx.Exec(`DROP INDEX chats_guest_id_index on chats;`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`alter table chats drop column expired_at;`)
			if err != nil {
				return err
			}
			_, err = tx.Exec(`alter table chats drop column owner;`)
			if err != nil {
				return err
			}
			_, err = tx.Exec(`alter table chats drop column guest_id;`)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
