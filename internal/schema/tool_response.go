package schema

type ToolListResponse struct {
	//Tool Tool `json:"tool"`
}

type ToolRemoteRequest struct {
	ApiKey       string         `json:"-"`
	FunctionName string         `json:"function_name"`
	Parameters   string         `json:"parameters"`
	User         *UserTokenInfo `json:"user"`
}

type ToolRemoteResponse struct {
	Success bool   `json:"success"`
	Content string `json:"content"`
}
