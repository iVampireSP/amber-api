package migrations

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

type chat struct {
	Id   int64  `json:"id"`
	Name string `xorm:"varchar(255) notnull" json:"name"`
}

func init() {
	migrations = append(migrations, &migrate.Migration{
		ID: "8",
		Migrate: func(tx *xorm.Engine) error {
			var err error
			var sql = `CREATE TABLE files (
  id   bigint unsigned AUTO_RANDOM,
  url varchar(255) DEFAULT NULL,
  url_hash varchar(255) DEFAULT NULL,
  file_hash varchar(255) DEFAULT NULL,
  mime_type varchar(255) DEFAULT NULL,
  path varchar(255) DEFAULT NULL,
  expired_at timestamp NULL DEFAULT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
			_, err = tx.Exec(sql)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index files_url_hash_index on files (url_hash);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index files_file_hash_index on files (file_hash);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index files_mime_type_index on files (mime_type);`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`create index files_expired_at_index on files (expired_at);`)
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			_, err := tx.Exec(`DROP TABLE files;`)
			if err != nil {
				return err
			}

			return nil
		},
	})
}
