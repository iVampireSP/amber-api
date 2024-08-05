package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "7",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
ALTER TABLE "public"."chats"  ADD COLUMN "expired_at"  timestamp(0);
ALTER TABLE "public"."chats" ADD COLUMN "owner" varchar(255) COLLATE "pg_catalog"."default";
ALTER TABLE "public"."chats" ADD COLUMN "guest_id" varchar(255) COLLATE "pg_catalog"."default";

-- set user_id null able
ALTER TABLE "public"."chats" ALTER COLUMN "user_id" DROP NOT NULL;

CREATE INDEX "chats_expired_at_index" ON "public"."chats" USING btree (
  "expired_at" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "chats_owner_index" ON "public"."chats" USING btree (
  "owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "chats_guest_id_index" ON "public"."chats" USING btree (
  "guest_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
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
