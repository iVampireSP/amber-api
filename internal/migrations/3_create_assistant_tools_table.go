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
CREATE TABLE "public"."assistant_tools" (
  "id" bigserial PRIMARY KEY,
  "assistant_id" bigint NOT NULL,
  "tool_id" bigint NOT NULL,
  "created_at" timestamp(0),
  "updated_at" timestamp(0)
);

CREATE INDEX "assistant_tools_assistant_id_tool_id_index" ON "public"."assistant_tools" USING btree (
  "assistant_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "tool_id" "pg_catalog"."int8_ops" ASC NULLS LAST
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
			_, err := tx.Exec("DROP TABLE IF EXISTS assistant_tools;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
