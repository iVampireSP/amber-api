package schema

type ToolListResponse struct {
	//Tool Tool `json:"tool"`
}

type ToolRemoteRequest struct {
	FunctionName string           `json:"function_name"`
	Parameters   interface{}      `json:"parameters"`
	User         *UserPublicInfo  `json:"user"`
	Chat         *ChatPublicModel `json:"chat"`
}

type ToolRemoteResponse struct {
	Success        bool `json:"success"`
	StopGeneration bool `json:"stop_generation"`
	// 记住响应。如果记住了响应，LLM 将会知道上一次工具的输出
	RememberResponse bool   `json:"remember_response"`
	Content          string `json:"content"`
}
