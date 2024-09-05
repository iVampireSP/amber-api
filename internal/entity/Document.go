package entity

import "rag-new/internal/schema"

type Document struct {
	Model
	Name      string           `json:"name"`
	Chunked   bool             `json:"chunked"`
	LibraryId schema.EntityId  `json:"library_id"`
	Library   *Library         `json:"library"`
	FileId    *schema.EntityId `json:"file_id"`
	File      *File            `json:"file"`
	// FileHash 是 File 结构体的 hash，用于判断文件内容是否发生变化
	// 只不过一般情况也不会改变，因为 File 就不会变
	FileHash *string `json:"file_hash"`
}

func (*Document) TableName() string {
	return "documents"
}
