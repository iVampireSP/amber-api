package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "2",
		Migrate: func(tx *xorm.Engine) error {
			var rawSQL = `
CREATE TABLE "public"."assistants" (
  "id" bigserial PRIMARY KEY ,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "prompt" text COLLATE "pg_catalog"."default" NOT NULL,
  "user_id" bigint NOT NULL,
  "created_at" timestamp(0),
  "updated_at" timestamp(0)
);

CREATE INDEX "assistants_user_id_index" ON "public"."assistants" USING btree (
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
			_, err := tx.Exec("DROP TABLE IF EXISTS assistants;")
			if err != nil {
				return err
			}
			return nil
		},
	})
}
