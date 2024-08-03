package schema

import (
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
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"  validate:"required"`
	HomepageUrl string `json:"homepage_url" validate:"url"`
	CallbackUrl string `json:"callback_url" validate:"url"`
	ToolId      int64  `json:"-"`
	Functions   []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  struct {
			Type       string `json:"type,omitempty"`
			Properties struct {
				Location struct {
					Type        string `json:"type"`
					Description string `json:"description"`
				} `json:"location"`
			} `json:"properties,omitempty"`
		} `json:"parameters"  validate:"required min=1"`
		Required []string `json:"required,omitempty" validate:"required"`
	} `json:"functions"`
}

func (td *ToolDiscoveryInput) Output() *ToolDiscoveryOutput {
	var output = ToolDiscoveryOutput{}
	var outputFunctions []ToolDiscoveryOutputFunctions

	output.Name = td.Name
	output.Description = td.Description
	output.HomepageUrl = td.HomepageUrl
	output.CallbackUrl = td.CallbackUrl

	// foreach
	for _, v := range td.Functions {
		var requires = make([]string, 0)

		if len(v.Required) > 0 {
			requires = v.Required
		}

		outputFunctions = append(outputFunctions, ToolDiscoveryOutputFunctions{
			Type: "function",
			Function: []ToolDiscoveryOutputFunction{
				{
					Name:        strconv.Itoa(int(td.ToolId)) + "_" + v.Name,
					Description: v.Description,
					Parameters:  v.Parameters,
					Required:    requires,
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
	Type     string `json:"type"`
	Function []ToolDiscoveryOutputFunction
}
type ToolDiscoveryOutputFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
	Required    []string    `json:"required"`
}

func (td *ToolDiscoveryOutput) FromDB(data []byte) error {
	return sonic.Unmarshal(data, &td)
}

func (td *ToolDiscoveryOutput) ToDB() ([]byte, error) {
	return sonic.Marshal(&td)
}
