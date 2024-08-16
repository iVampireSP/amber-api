package schema

import (
	"github.com/bytedance/sonic"
	"github.com/tmc/langchaingo/llms"
)

type ResponseState string

const (
	StateChunk        ResponseState = "chunk"
	StateToolCalling  ResponseState = "tool_calling"
	StateToolResponse ResponseState = "tool_response"
	StateToolSuccess  ResponseState = "tool_success"
	StateToolFailed   ResponseState = "tool_failed"
	StateFinished     ResponseState = "finished" // finished 为全部完成
	StateDone         ResponseState = "done"     // done 为一轮请求完成
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
	TokenUsage          *TokenUsage          `json:"token_usage"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type FunctionChunk struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type FunctionCallArgs map[string]interface{}

func (f *FunctionCallArgs) JSON() ([]byte, error) {
	return sonic.Marshal(f)
}

type LLMChat struct {
	ResponseChan   chan *AssistantResponse
	SystemPrompt   string
	UserPublicInfo *UserPublicInfo
	Tools          []llms.Tool
	MaxTokens      int     `json:"max_tokens,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
	TopK           int     `json:"top_k,omitempty"`
	TopP           float64 `json:"top_p,omitempty"`
	N              int     `json:"n,omitempty"`
}
