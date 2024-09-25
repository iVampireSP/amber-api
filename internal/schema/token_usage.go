package schema

type SiteUsageResponse struct {
	TodayTokens        int `json:"today_tokens"`
	TodayToolCalls     int `json:"today_tool_calls"`
	YesterdayTokens    int `json:"yesterday_tokens"`
	YesterdayToolCalls int `json:"yesterday_tool_calls"`
}
