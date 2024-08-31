package schema

import (
	"database/sql/driver"
	"github.com/bytedance/sonic"
	"strconv"
)

type ToolCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	Description string `json:"description" binding:"required" validate:"max=255"`
	Url         string `json:"url" binding:"required" validate:"max=255"`
	ApiKey      string `json:"api_key" validate:"max=255"`
}

type ToolUpdateRequest struct {
}

type ToolDiscoveryInput struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"  validate:"required"`
	HomepageUrl string            `json:"homepage_url" validate:"url"`
	CallbackUrl string            `json:"callback_url" validate:"url"`
	ToolId      EntityId          `json:"-"`
	Functions   []*FunctionsInput `json:"functions"`
}

type FunctionsInput struct {
	Name        string                                `json:"name" validate:"required"`
	Description string                                `json:"description" validate:"required"`
	Parameters  ToolDiscoveryOutputFunctionParameters `json:"parameters" validate:"required"`
}

func (td *ToolDiscoveryInput) Output() *ToolDiscoveryOutput {
	var output = ToolDiscoveryOutput{}
	var outputFunctions []ToolDiscoveryOutputFunctions

	output.Name = td.Name
	output.Description = td.Description
	output.HomepageUrl = td.HomepageUrl
	output.CallbackUrl = td.CallbackUrl

	for _, v := range td.Functions {
		// 容忍处理，OpenAI 官方不允许 required 为 null 或缺失的情况。
		// 本程序如果检测到了这种情况，将 required 设置为 []
		var requires = make([]string, 0)

		if len(v.Parameters.Required) == 0 {
			v.Parameters.Required = requires
		}

		outputFunctions = append(outputFunctions, ToolDiscoveryOutputFunctions{
			Type: "function",
			Functions: []*ToolDiscoveryOutputFunction{
				{
					Name:        strconv.Itoa(int(td.ToolId)) + "_" + v.Name,
					Description: v.Description,
					Parameters:  v.Parameters,
				},
			},
		})
	}

	output.ToolFunctions = outputFunctions

	return &output
}

type ToolDiscoveryOutput struct {
	Name          string                         `json:"name" `
	HomepageUrl   string                         `json:"homepage_url" `
	CallbackUrl   string                         `json:"callback_url"`
	Description   string                         `json:"description"`
	ToolFunctions []ToolDiscoveryOutputFunctions `json:"function"`
}
type ToolDiscoveryOutputFunctions struct {
	Type      string                         `json:"type"`
	Functions []*ToolDiscoveryOutputFunction `json:"functions"`
}
type ToolDiscoveryOutputFunction struct {
	Name        string                                `json:"name"`
	Description string                                `json:"description"`
	Parameters  ToolDiscoveryOutputFunctionParameters `json:"parameters"`
}

type ToolDiscoveryOutputFunctionParameters struct {
	Type       string      `json:"type,omitempty" validate:"required"`
	Properties interface{} `json:"properties" validate:"required"`
	Required   []string    `json:"required" validate:"required"`
}

func (td *ToolDiscoveryOutput) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &td)
}

func (td ToolDiscoveryOutput) Value() (driver.Value, error) {
	return sonic.Marshal(&td)
}
