package orm

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/yxlimo/xormzap"
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
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Name)

	engine, err := xorm.NewEngine("postgres", dsn)
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
