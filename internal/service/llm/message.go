package llm

import "github.com/tmc/langchaingo/llms"

type Message struct {
	HasFile        bool
	MessageContent []llms.MessageContent
}
