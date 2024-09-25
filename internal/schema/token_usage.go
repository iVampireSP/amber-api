package schema

type SiteUsageResponse struct {
	MonthTokens    int `json:"month_tokens"`
	MonthToolCalls int `json:"month_tool_calls"`
}
