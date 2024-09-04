package message

type Type string

const (
	Chunk       Type = "chunk"
	MemoryAdded Type = "memory_added"
)

func (t Type) String() string {
	return string(t)
}
