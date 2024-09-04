package openai

import (
	"github.com/sashabaranov/go-openai"
	"rag-new/internal/base/conf"
)

func NewOpenAI(c *conf.Config) *openai.Client {
	config := openai.DefaultConfig(c.OpenAI.ApiKey)
	config.BaseURL = c.OpenAI.BaseUrl

	return openai.NewClientWithConfig(config)
}
