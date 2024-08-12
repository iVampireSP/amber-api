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
ALTER TABLE "public"."chat_messages" RENAME COLUMN "input_tokens" TO "prompt_tokens";
ALTER TABLE "public"."chat_messages" RENAME COLUMN "output_tokens" TO "completion_tokens";
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
-- Drop indexes
DROP INDEX "chats_expired_at_index";
DROP INDEX "chats_owner_index";
DROP INDEX "chats_guest_id_index";

-- Drop columns
ALTER TABLE "public"."chats" DROP COLUMN "expired_at";
ALTER TABLE "public"."chats" DROP COLUMN "owner";
ALTER TABLE "public"."chats" DROP COLUMN "guest_id";
`)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
