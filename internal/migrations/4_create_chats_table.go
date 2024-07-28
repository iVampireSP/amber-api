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
CREATE TABLE "public"."chats" (
  "id" bigserial PRIMARY KEY ,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "assistant_id" int8 NOT NULL,
  "user_id" int8 NOT NULL,
  "created_at" timestamp(0),
  "updated_at" timestamp(0)
);

CREATE INDEX "chats_assistant_id_user_id_index" ON "public"."chats" USING btree (
  "assistant_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
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
			_, err := tx.Exec("DROP TABLE IF EXISTS chats;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
