package conf

import (
	"os"
	"rag-new/configs"
)

type Config struct {
	Http *Http `yaml:"http"`

	Debug *Debug `yaml:"debug"`

	Database *Database `yaml:"database"`

	Redis *Redis `yaml:"redis"`

	JWKS *JWKS `yaml:"jwks"`

	Metrics *Metrics `yaml:"metrics"`

	OpenAI *OpenAI `yaml:"openai"`

	LLM *LLM `yaml:"llm"`

	S3 *S3 `yaml:"s3"`
}

type Http struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Url  string `yaml:"url"`
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
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWKS struct {
	Url string `yaml:"url"`
}

type Metrics struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
}

type OpenAI struct {
	ApiKey      string `yaml:"api_key" mapstructure:"api_key"`
	BaseUrl     string `yaml:"api_base" mapstructure:"base_url"`
	Model       string `yaml:"model" mapstructure:"model"`
	VisionModel string `yaml:"vision_model" mapstructure:"vision_model"`
}

type S3 struct {
	Endpoint         string `yaml:"endpoint" mapstructure:"endpoint"`
	ExternalEndpoint string `yaml:"external_endpoint" mapstructure:"external_endpoint"`
	AccessKey        string `yaml:"access_key" mapstructure:"access_key"`
	SecretKey        string `yaml:"secret_key" mapstructure:"secret_key"`
	Bucket           string `yaml:"bucket" mapstructure:"bucket"`
	UseSSL           bool   `yaml:"use_ssl" mapstructure:"use_ssl"`
	Region           string `yaml:"region" mapstructure:"region"`
}

type LLM struct {
	MaxTokens int `yaml:"max_tokens"`
	//Temperature float64 `yaml:"temperature"`
	//TopP float64 `yaml:"top_p"`
	//N int `yaml:"n"`
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
