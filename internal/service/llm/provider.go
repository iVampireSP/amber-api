package llm

import (
	"github.com/tmc/langchaingo/llms/openai"
	"rag-new/internal/base/conf"
)

type FunctionChunk struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type FunctionCallArgs map[string]interface{}

type ResponseState string

const (
	StateChunk        ResponseState = "chunk"
	StateToolCalling  ResponseState = "tool_calling"
	StateToolResponse ResponseState = "tool_response"
	StateToolSuccess  ResponseState = "tool_success"
	StateToolFailed   ResponseState = "tool_failed"
	StateToolCalled   ResponseState = "tool_called"
	StateDone         ResponseState = "done"
	StateFailed       ResponseState = "failed"
)

type ChunkMessage struct {
	Content string `json:"content"`
}

type ToolCallMessage struct {
	Name string           `json:"name"`
	Args FunctionCallArgs `json:"args"`
}

type ToolResponseMessage struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type AssistantResponse struct {
	State               ResponseState        `json:"state"`
	ChunkMessage        *ChunkMessage        `json:"chunk_message"`
	ToolCallMessage     *ToolCallMessage     `json:"tool_call_message"`
	ToolResponseMessage *ToolResponseMessage `json:"tool_response_message"`
	Content             string               `json:"content"`
}

type Service struct {
	OpenAI *openai.LLM
}

func NewLLM(config *conf.Config) *Service {
	llm, err := openai.New(
		openai.WithToken(config.OpenAI.ApiKey),
		openai.WithBaseURL(config.OpenAI.BaseUrl),
		openai.WithModel(config.OpenAI.Model),
	)

	if err != nil {
		panic(err)
	}

	return &Service{llm}
}
