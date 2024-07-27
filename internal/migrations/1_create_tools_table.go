package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "1",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE "public"."tools" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "discovery_url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "api_key" varchar(255) COLLATE "pg_catalog"."default",
  "data" json,
  "user_id" BIGINT NOT NULL,
  "created_at" timestamp(0),
  "updated_at" timestamp(0)
);

CREATE INDEX "tools_discovery_url_index" ON "public"."tools" USING btree (
  "discovery_url" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "tools_user_id_index" ON "public"."tools" USING btree (
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
			_, err := tx.Exec("DROP TABLE IF EXISTS tools;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
