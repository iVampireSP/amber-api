package orm

import (
	"fmt"
	"github.com/yxlimo/xormzap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"xorm.io/xorm"
)

//
//func ProviderOrm() *Orm {
//	return NewOrm()
//}

func NewXORM(
	config *conf.Config,
	logger *logger.Logger,
) (*xorm.Engine, error) {
	//dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}

	engine.SetLogger(xormzap.Logger(logger.Logger))

	engine.ShowSQL(config.Debug.Enabled)

	err = engine.Ping()
	if err != nil {
		panic(err)
	}

	return engine, nil
}

func NewGORM(
	config *conf.Config,
	logger *logger.Logger,
) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)

	zapGormLogger := zapgorm2.New(logger.Logger)
	zapGormLogger.SetAsDefault()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: zapGormLogger,
	})

	if err != nil {
		panic(err)
	}

	return db
}
