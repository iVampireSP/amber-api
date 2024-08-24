package builtin_tool

import (
	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
)

var tools = []llms.Tool{
	//{
	//	Type: "function",
	//	Function: &llms.FunctionDefinition{
	//		Name:        prefix("now"),
	//		Description: "get current time",
	//		Parameters: jsonschema.Definition{
	//			Type: jsonschema.Object,
	//			Properties: map[string]jsonschema.Definition{
	//				"timezone": {
	//					Type:        jsonschema.String,
	//					Description: "Timezone, default Asia/Shanghai",
	//				},
	//			},
	//			// 必须要为 1
	//			Required: make([]string, 1),
	//		},
	//	},
	//},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("describe_image"),
			Description: "describe image by natural language",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"image_id": {
						Type:        jsonschema.Integer,
						Description: "The id of the image you want to describe, must get from history.",
					},
					"question": {
						Type:        jsonschema.String,
						Description: "What you need to explain, using natural language like 'What is this image about?', Write questions in the user language like '我想知道这张图片是什么关于的'。",
					},
				},
				Required: []string{
					"image_id",
					"question",
				},
			},
		},
	},
}
