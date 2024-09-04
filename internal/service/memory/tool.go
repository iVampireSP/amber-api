package memory

import (
	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
)

var tools = []llms.Tool{
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "add_memory",
			Description: "Add a memory",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"data": {
						Type:        jsonschema.String,
						Description: "Data to add to memory",
					},
				},
				Required: []string{
					"data",
				},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "update_memory",
			Description: "Update memory provided ID and data",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"memory_id": {
						Type:        jsonschema.Integer,
						Description: "memory_id of the memory to update",
					},
					"data": {
						Type:        jsonschema.String,
						Description: "Updated data for the memory",
					},
				},
				Required: []string{
					"memory_id",
					"data",
				},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "delete_memory",
			Description: "Delete memory by memory_id",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"memory_id": {
						Type:        jsonschema.Integer,
						Description: "memory_id of the memory to delete",
					},
				},
				Required: []string{
					"memory_id",
				},
			},
		},
	},
}
