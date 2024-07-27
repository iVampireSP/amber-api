package conf

import (
	"os"
	"rag-new/configs"
)

type Config struct {
	Http *Http `yaml:"http"`

	Debug *Debug `yaml:"debug"`

	Database *Database `yaml:"database"`
}

type Http struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Debug struct {
	Enabled bool `yaml:"enabled"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SslMode  string `yaml:"sslmode"`
}

func CreateConfigIfNotExists() {
	// create if not exists
	var configName = "config.yaml"

	if _, err := os.Stat(configName); os.IsNotExist(err) {
		f, err := os.Create(configName)
		if err != nil {
			panic(err)
		}

		// write default from embed
		_, err = f.Write(configs.Config)
		if err != nil {
			panic(err)
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)
	}
}
