package conf

import (
	"github.com/spf13/viper"
	"rag-new/internal/base/logger"
)

func ProviderConfig(logger *logger.Logger) *Config {
	CreateConfigIfNotExists()
	c := &Config{}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		logger.Sugar.Fatal(err)
	}
	if err := v.Unmarshal(c); err != nil {
		panic(err)
	}

	return c
}
