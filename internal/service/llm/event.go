package llm

import (
	"context"
	"rag-new/internal/message"
	"rag-new/internal/schema"
)

func (s *Service) event(ctx context.Context, llmChat *schema.LLMChat) {
	msg, ok := <-llmChat.ResponseChan
	if !ok {
		return
	}

	if msg == nil {
		return
	}

	switch msg.State {
	case schema.StateChunk:
		chunkMessage := s.message.NewChunk(msg.Content, llmChat.UserPublicInfo)
		err := s.streamService.SendEvent(ctx, message.Chunk.String(), chunkMessage)
		if err != nil {
			s.Logger.Sugar.Errorf("Unable send event: %v", err)
			return
		}
		return
	case schema.StateToolSuccess:
		return
	case schema.StateToolCalling:

		return
	case schema.StateToolResponse:

		return
	case schema.StateDone:
		return
	case schema.StateFailed:
		return
	case schema.StateFinished:
		return
	case schema.StateToolFailed:
		return
	default:
		return
	}
}
