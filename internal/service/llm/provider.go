package llm

import (
	"github.com/bytedance/sonic"
	"github.com/tmc/langchaingo/llms/openai"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/tool"
)

type FunctionChunk struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type FunctionCallArgs map[string]interface{}

func (f *FunctionCallArgs) JSON() ([]byte, error) {
	return sonic.Marshal(f)
}

type ResponseState string

const (
	StateChunk        ResponseState = "chunk"
	StateToolCalling  ResponseState = "tool_calling"
	StateToolResponse ResponseState = "tool_response"
	StateToolSuccess  ResponseState = "tool_success"
	StateToolFailed   ResponseState = "tool_failed"
	StateToolCalled   ResponseState = "tool_called"
	StateFinished     ResponseState = "finished"
	StateFailed       ResponseState = "failed"
)

type ChunkMessage struct {
	Content string `json:"content"`
}

type ToolCallMessage struct {
	ToolName     string           `json:"tool_name"`
	FunctionName string           `json:"function_name"`
	Args         FunctionCallArgs `json:"args"`
}

type ToolResponseMessage struct {
	ToolName     string `json:"tool_name"`
	FunctionName string `json:"function_name"`
	Content      string `json:"content"`
}

type AssistantResponse struct {
	State               ResponseState        `json:"state"`
	ChunkMessage        *ChunkMessage        `json:"chunk_message"`
	ToolCallMessage     *ToolCallMessage     `json:"tool_call_message"`
	ToolResponseMessage *ToolResponseMessage `json:"tool_response_message"`
	Content             string               `json:"content"`
}

type Service struct {
	OpenAI           *openai.LLM
	Logger           *logger.Logger
	AssistantService *assistant.Service
	ToolService      *tool.Service
}

func NewLLM(config *conf.Config, logger *logger.Logger, assistantService *assistant.Service, toolService *tool.Service) *Service {
	llm, err := openai.New(
		openai.WithToken(config.OpenAI.ApiKey),
		openai.WithBaseURL(config.OpenAI.BaseUrl),
		openai.WithModel(config.OpenAI.Model),
	)

	if err != nil {
		panic(err)
	}

	return &Service{llm, logger, assistantService, toolService}
}
