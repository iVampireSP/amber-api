package message

import (
	"github.com/bytedance/sonic"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"time"
)

type MemoryAddedMessage struct {
	Type      Type            `json:"type"`
	MemoryId  schema.EntityId `json:"memory_id"`
	Content   string          `json:"content"`
	Model     string          `json:"model"`
	UserId    schema.UserId   `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
}

func (m *Message) NewMemoryAdded(memory *entity.Memory) *MemoryAddedMessage {
	return &MemoryAddedMessage{
		Type:      MemoryAdded,
		MemoryId:  memory.Id,
		Content:   memory.Content,
		Model:     memory.EmbeddingModel,
		UserId:    memory.UserId,
		CreatedAt: memory.CreatedAt,
	}
}

func (m *MemoryAddedMessage) JSON() ([]byte, error) {
	return sonic.Marshal(m)
}
