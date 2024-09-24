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
			Name: prefix("describe_image"),
			Description: "only used to retrieve the content of images and cannot obtain file information of images." +
				" only for which mimetype is image",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"hash": {
						Type:        jsonschema.Integer,
						Description: "The hash of the file(with image mimetype, from history) you want to describe",
					},
					"url": {
						Type: jsonschema.String,
						Description: "The url of the image you want to describe" +
							"(URL or hash must be chosen between two options)",
					},
					"prompt": {
						Type:        jsonschema.String,
						Description: "What you need to explain.",
					},
				},
				Required: []string{
					"question",
				},
			},
		},
	},
	//{
	//	Type: "function",
	//	Function: &llms.FunctionDefinition{
	//		Name:        prefix("download_file"),
	//		Description: "download file from url",
	//		Parameters: jsonschema.Definition{
	//			Type: jsonschema.Object,
	//			Properties: map[string]jsonschema.Definition{
	//				"url": {
	//					Type: jsonschema.Integer,
	//					Description: "the url of the file you want to download. " +
	//						"when downloaded, you can get file id from history",
	//				},
	//			},
	//			Required: []string{
	//				"url",
	//			},
	//		},
	//	},
	//},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("search_web"),
			Description: "Search the internet",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"query": {
						Type:        jsonschema.String,
						Description: "the query you want to search",
					},
				},
				Required: []string{
					"query",
				},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("get_url_content"),
			Description: "Get the website content of the url",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"url": {
						Type:        jsonschema.String,
						Description: "the url of the website you want to get content",
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
			Description: "It's useful for generating/drawing images,if there are no special requirements, always use the markdown syntax to display images",
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
			Name: prefix("calculator"),
			Description: "It's useful for mathematical calculations," +
				"every time a step is executed, the user must be informed and then proceed to the next step." +
				"When encountering functions such as log and sqrt, " +
				"you need to call the tool multiple times to calculate, " +
				"the calculation process must be written out before calling the tool to perform the calculation." +
				"The result of this tool is always right.",
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
