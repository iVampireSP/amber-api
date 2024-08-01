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
CREATE TABLE "public"."chat_messages" (
  "id" bigserial PRIMARY KEY ,
  "chat_id" int8 NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "role" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "input_tokens" int4,
  "output_tokens" int4,
  "total_tokens" int4,
  "created_at" timestamp(0),
  "updated_at" timestamp(0),
  "tool_calls" json
);

CREATE INDEX "chat_messages_chat_id_role_index" ON "public"."chat_messages" USING btree (
  "chat_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "role" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
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
			_, err := tx.Exec("DROP TABLE IF EXISTS chat_messages;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
