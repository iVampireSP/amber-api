package memory

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
)

// executeToolCalls executes the tool calls in the response and returns the
// updated message history.
func (s *Service) executeToolCalls(ctx context.Context, userId schema.UserId, messageHistory []llms.MessageContent, resp *llms.ContentResponse) ([]llms.MessageContent, error) {
	for _, toolCall := range resp.Choices[0].ToolCalls {
		s.Logger.Sugar.Infof("Memory Calling: %s", toolCall.FunctionCall.Name)
		switch toolCall.FunctionCall.Name {
		case "add_memory":
			var args struct {
				Data string `json:"data"`
			}
			if err := sonic.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				return nil, err
			}

			memEntity, err := s.addMemory(ctx, args.Data, userId)
			if err != nil {
				return nil, err
			}

			s.Logger.Sugar.Infof("Memory added, id: %d, content: %s", memEntity.Id, memEntity.Content)

		case "update_memory":
			var args struct {
				Data     string          `json:"data"`
				MemoryId schema.EntityId `json:"memory_id"`
			}
			if err := sonic.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				return nil, err
			}

			_, err := s.updateMemory(ctx, args.MemoryId, args.Data)
			if err != nil {
				return nil, err
			}
		case "delete_memory":
			var args struct {
				MemoryId schema.EntityId `json:"memory_id"`
			}
			if err := sonic.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				return nil, err
			}

			err := s.deleteMemory(ctx, args.MemoryId)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported tool: %s", toolCall.FunctionCall.Name)
		}
	}

	return messageHistory, nil
}
