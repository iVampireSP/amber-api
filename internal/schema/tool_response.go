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
	Success        bool   `json:"success"`
	StopGeneration bool   `json:"stop_generation"`
	Content        string `json:"content"`
}
