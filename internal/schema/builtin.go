package schema

type CallBuiltInToolRequest struct {
	FunctionName string
	Args         FunctionCallArguments
}

type CallBuiltInResponse struct {
	Success          bool     `json:"success"`
	StopGeneration   bool     `json:"stop_generation"`
	RememberResponse bool     `json:"remember_response"`
	Content          string   `json:"content"`
	Append           bool     `json:"-"`
	Role             ChatRole `json:"-"`
	Text             string   `json:"-"`
	*TokenUsage
}
