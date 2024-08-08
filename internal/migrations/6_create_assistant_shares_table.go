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
CREATE TABLE "public"."assistant_shares" (
  "id" bigserial PRIMARY KEY ,
  "assistant_id" int8 NOT NULL,
  "token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL UNIQUE,
  "created_at" timestamp(0),
  "updated_at" timestamp(0)
);

CREATE INDEX "assistant_shares_assistant_id_token_index" ON "public"."assistant_shares" USING btree (
  "assistant_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "token" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST 
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
			_, err := tx.Exec("DROP TABLE IF EXISTS assistant_shares;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
