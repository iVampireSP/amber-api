package message

import (
	"github.com/bytedance/sonic"
	"rag-new/internal/schema"
)

type ChunkMessage struct {
	Type    Type                   `json:"type"`
	Content string                 `json:"content"`
	User    *schema.UserPublicInfo `json:"user"`
}

func (m *Message) NewChunk(content string, user *schema.UserPublicInfo) *ChunkMessage {
	return &ChunkMessage{
		Type:    Chunk,
		User:    user,
		Content: content,
	}
}

func (m *ChunkMessage) JSON() ([]byte, error) {
	return sonic.Marshal(m)
}
