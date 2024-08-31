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
	// 以下为隐藏字段，暂未开放给 API 调用的工具
	// 增加的聊天记录
	Append bool     `json:"-"`
	Role   ChatRole `json:"-"`
	Text   string   `json:"-"`
}
