package llm

import (
	"context"
	"rag-new/internal/message"
	"rag-new/internal/schema"
)

func (s *Service) event(ctx context.Context, user *schema.UserPublicInfo, msg *schema.AssistantResponse) {
	switch msg.State {
	case schema.StateChunk:
		// 这个会顺序错误
		// 但是目前没关系，因为还没到真正的用途
		chunkMessage := s.message.NewChunk(msg.Content, user)
		err := s.streamService.SendEvent(ctx, message.Chunk.String(), chunkMessage)
		if err != nil {
			s.Logger.Sugar.Errorf("Unable send event: %v", err)
		}

		return
		//case schema.StateToolSuccess:
		//	return
		//case schema.StateToolCalling:
		//	return
		//case schema.StateToolResponse:
		//	return
		//case schema.StateDone:
		//	return
		//case schema.StateFailed:
		//	return
		//case schema.StateFinished:
		//	return
		//case schema.StateToolFailed:
		//	return
		//default:
		//	return
	}
}
