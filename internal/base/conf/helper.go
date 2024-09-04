package conf

import (
	"os"
	"rag-new/configs"
	"rag-new/pkg/consts"
)

func (o *OpenAI) CanUse(model string) bool {
	if o.AllowedModels == nil {
		return true
	}

	if model == consts.AutoModel {
		return true
	}

	for _, allowedModel := range o.AllowedModels {
		if allowedModel == model {
			return true
		}
	}

	return false
}

func createConfigIfNotExists(path string) {
	if path != "" {
		return
	}

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
