package llm

import "github.com/tmc/langchaingo/llms"

type Message struct {
	HasImage       bool
	MessageContent []llms.MessageContent
}
