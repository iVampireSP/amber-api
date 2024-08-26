package schema

type OpenAIChatCompletionRequest struct {
	// Required
	Model string `json:"model"`
	// Required
	Messages []OpenAIChatCompletionRequestMessage `json:"messages"`
	// Optional
	MaxTokens int `json:"max_tokens,omitempty"`
	// Optional
	Temperature float64 `json:"temperature,omitempty"`
	// Optional
	TopP float64 `json:"top_p,omitempty"`
	// Optional
	N int `json:"n,omitempty"`
	// Optional
	Stream bool `json:"stream,omitempty"`
}

type OpenAIChatCompletionRequestMessage struct {
	// Required
	Role string `json:"role"`
	// Required
	Content  any             `json:"content"`
	ImageUrl *OpenAIImageUrl `json:"image_url"`
}

type OpenAIImageUrl struct {
	Url    string `json:"url"`
	Detail string `json:"detail"`
}

type OpenAIChatCompletionResponseChoice struct {
	Index   int32                              `json:"index"`
	Message OpenAIChatCompletionRequestMessage `json:"message"`
}

type OpenAIChatCompletionStreamDelta struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type OpenAIChatCompletionStreamResponseChoice struct {
	Delta        OpenAIChatCompletionStreamDelta `json:"delta"`
	Index        int                             `json:"index"`
	FinishReason interface{}                     `json:"finish_reason"`
}

type OpenAIChatCompletionResponse struct {
	ID      string                               `json:"id"`
	Object  string                               `json:"object"`
	Created int64                                `json:"created"`
	Choices []OpenAIChatCompletionResponseChoice `json:"choices"`
	Usage   *TokenUsage                          `json:"usage"`
	Model   string                               `json:"model"`
}

type OpenAIChatCompletionStreamResponse struct {
	ID      string                                     `json:"id"`
	Object  string                                     `json:"object"`
	Created int64                                      `json:"created"`
	Choices []OpenAIChatCompletionStreamResponseChoice `json:"choices"`
	Usage   *TokenUsage                                `json:"usage"`
	Model   string                                     `json:"model"`
}
