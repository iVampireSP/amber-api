package conf

import (
	"github.com/spf13/viper"
	"log"
)

func ProviderConfig() *Config {
	c := &Config{}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	if err := v.Unmarshal(c); err != nil {
		panic(err)
	}

	return c
}
