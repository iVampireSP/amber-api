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
			Description: "only used to retrieve the content of images and cannot obtain file information of images. only for which mimetype is image",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"file_id": {
						Type:        jsonschema.Integer,
						Description: "The id of the file(with image mimetype, from history) you want to describe",
					},
					"url": {
						Type:        jsonschema.String,
						Description: "The url of the image you want to describe(URL and file ID must be chosen between two options)",
					},
					"question": {
						Type:        jsonschema.String,
						Description: "What you need to explain, using natural language like 'What is this image about?', Write questions in the user language like '我想知道这张图片是什么关于的'。",
					},
				},
				Required: []string{
					"question",
				},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("download_file"),
			Description: "download file from url",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"url": {
						Type:        jsonschema.Integer,
						Description: "the url of the file you want to download. when downloaded, you can get file id from history",
					},
				},
				Required: []string{
					"url",
				},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("generate_image"),
			Description: "It's useful for generating/drawing images",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"prompt": {
						Type:        jsonschema.String,
						Description: "prompt to generate image",
					},
					"size": {
						Type:        jsonschema.String,
						Description: "size of image",
						Enum:        dallEAllowedSizes,
					},
				},
				Required: []string{
					"prompt", "size",
				},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("calculator"),
			Description: "It's useful for mathematical calculations",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"A": {
						Type:        jsonschema.String,
						Description: "Number A",
					},
					"B": {
						Type:        jsonschema.String,
						Description: "Number B",
					},
					"Method": {
						Type:        jsonschema.String,
						Description: "Method",
						Enum:        calculatorAllowedMethods,
					},
				},
				Required: []string{
					"A", "B", "Method",
				},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("compare"),
			Description: "It's useful for comparing numbers",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"A": {
						Type:        jsonschema.String,
						Description: "Number A",
					},
					"B": {
						Type:        jsonschema.String,
						Description: "Number B",
					},
				},
				Required: []string{
					"A", "B",
				},
			},
		},
	},
}
