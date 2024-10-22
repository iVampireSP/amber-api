package builtin_tool

import (
	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
)

type WithoutOptions struct {
	Image    bool
	Browsing bool
}

func (s *Service) GetTools(without *WithoutOptions) []llms.Tool {
	var tools []llms.Tool

	// 如果不包含图片工具
	//if !without.Image {
	//	tools = append(tools, llms.Tool{
	//		Type: "function",
	//		Function: &llms.FunctionDefinition{
	//			Name: prefix("describe_image"),
	//			Description: "only used to retrieve the content of images and cannot obtain file information of images." +
	//				" only for which mimetype is image",
	//			Parameters: jsonschema.Definition{
	//				Type: jsonschema.Object,
	//				Properties: map[string]jsonschema.Definition{
	//					"hash": {
	//						Type:        jsonschema.String,
	//						Description: "The hash of the file(with image mimetype, from history) you want to describe",
	//					},
	//					"url": {
	//						Type: jsonschema.String,
	//						Description: "The url of the image you want to describe" +
	//							"(URL or hash must be chosen between two options)",
	//					},
	//					"prompt": {
	//						Type:        jsonschema.String,
	//						Description: "What you need to explain.",
	//					},
	//				},
	//				Required: []string{
	//					"prompt",
	//				},
	//			},
	//		},
	//	})
	//}

	tools = append(tools, llms.Tool{
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
	})

	//tools = append(tools, llms.Tool{
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
	//})

	// 如果禁用网页浏览
	if !without.Browsing {
		tools = append(tools, llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        prefix("browser"),
				Description: "Browsing the internet",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"query_or_url": {
							Type:        jsonschema.String,
							Description: "Provide a search keyword or a webpage URL. If a URL is provided, the content of the webpage will be retrieved; otherwise, a search engine will be used.",
						},
					},
					Required: []string{
						"query_or_url",
					},
				},
			},
		})
	}

	// 如果禁用网页浏览
	// 这个应该废弃，我们应该将网页浏览和查询放在一个工具里面
	//if !without.Browsing {
	//	tools = append(tools, llms.Tool{
	//		Type: "function",
	//		Function: &llms.FunctionDefinition{
	//			Name:        prefix("browser_url"),
	//			Description: "Browser the web",
	//			Parameters: jsonschema.Definition{
	//				Type: jsonschema.Object,
	//				Properties: map[string]jsonschema.Definition{
	//					"url": {
	//						Type:        jsonschema.String,
	//						Description: "the url of the website you want to get content",
	//					},
	//				},
	//				Required: []string{
	//					"url",
	//				},
	//			},
	//		},
	//	})
	//}

	tools = append(tools, llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        prefix("calculator"),
			Description: "Useful for getting the result of a math expression or comparing numbers. The input to this tool should be a valid mathematical expression that could be executed by a simple calculator.",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"number_a": {
						Type:        jsonschema.String,
						Description: "Number A",
					},
					"operator": {
						Type:        jsonschema.String,
						Description: "Operator",
						Enum:        calculatorAllowedMethods,
					},
					"number_b": {
						Type:        jsonschema.String,
						Description: "Number B",
					},
				},
				Required: []string{
					"number_a", "operator", "number_b",
				},
			},
		},
	})

	//tools = append(tools, llms.Tool{
	//	Type: "function",
	//	Function: &llms.FunctionDefinition{
	//		Name:        prefix("compare"),
	//		Description: "It's useful for comparing numbers",
	//		Parameters: jsonschema.Definition{
	//			Type: jsonschema.Object,
	//			Properties: map[string]jsonschema.Definition{
	//				"number_a": {
	//					Type:        jsonschema.String,
	//					Description: "Number A",
	//				},
	//				"number_b": {
	//					Type:        jsonschema.String,
	//					Description: "Number B",
	//				},
	//			},
	//			Required: []string{
	//				"number_a", "number_b",
	//			},
	//		},
	//	},
	//})

	return tools
}
