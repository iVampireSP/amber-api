package schema

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/mitchellh/mapstructure"
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
	StateMessage      ResponseState = "message" // message 为一条消息
)

type ChunkMessage struct {
	Content string `json:"content"`
}

type ToolCallMessage struct {
	ToolName     string                `json:"tool_name"`
	FunctionName string                `json:"function_name"`
	Arguments    FunctionCallArguments `json:"args"`
}

type ToolResponseMessage struct {
	ToolName       string                `json:"tool_name"`
	FunctionName   string                `json:"function_name"`
	Arguments      FunctionCallArguments `json:"arguments"`
	StopGeneration bool                  `json:"stop_generation"`
	Content        string                `json:"content"`
	Append         bool                  `json:"-"`
	Role           ChatRole              `json:"-"`
	Text           string                `json:"-"`
}

type AssistantResponse struct {
	State               ResponseState        `json:"state"`
	ChunkMessage        *ChunkMessage        `json:"chunk_message"`
	ToolCallMessage     *ToolCallMessage     `json:"tool_call_message"`
	ToolResponseMessage *ToolResponseMessage `json:"tool_response_message"`
	Content             string               `json:"content"`
	TokenUsage          *TokenUsage          `json:"token_usage"`
	Internal            *AssistantInternal   `json:"-"`
}

type AssistantInternal struct {
	ToolCall   *llms.ToolCall
	ToolCallId string
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

type FunctionCallArguments map[string]interface{}

func (f *FunctionCallArguments) JSON() ([]byte, error) {
	return sonic.Marshal(f)
}

func (f *FunctionCallArguments) Unmarshal(out interface{}) error {
	return mapstructure.Decode(f, out)
}

func (f *FunctionCallArguments) String() (string, error) {
	j, err := f.JSON()

	return string(j), err
}

type LLMChat struct {
	ResponseChan   chan *AssistantResponse
	SystemPrompt   string
	UserPublicInfo *UserPublicInfo
	Chat           *ChatPublicModel `json:"chat"`
	//ToolCallToken  string           `json:"tool_call_token"`
	Tools       []llms.Tool
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	N           int     `json:"n,omitempty"`
	Model       string  `json:"model"`
	//WithoutImage    bool    `json:"-"`
	WithoutBrowsing bool `json:"-"`
	//UseVisionModel  bool `json:"-"`
}

type ChatPublicModel struct {
	Name        string     `json:"name"`
	ID          EntityId   `json:"id"`
	AssistantId *EntityId  `json:"assistant_id"`
	ExpiredAt   *time.Time `json:"expired_at"`
	Owner       ChatOwner  `json:"owner"`
	GuestId     *string    `json:"guest_id"`
}
