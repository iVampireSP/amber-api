package orm

import (
	"fmt"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

//
//func ProviderOrm() *Orm {
//	return NewOrm()
//}

func NewGORM(
	config *conf.Config,
	logger *logger.Logger,
) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name)

	if config.Debug.Enabled {
		db, err := gorm.Open(postgres.Open(dsn))

		if err != nil {
			panic(err)
		}

		return db
	}

	zapGormLogger := zapgorm2.New(logger.Logger)
	zapGormLogger.SetAsDefault()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: zapGormLogger,
	})

	if err != nil {
		panic(err)
	}

	return db
}
