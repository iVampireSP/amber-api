package schema

type FileDownloadRequest struct {
	FileId int64 `uri:"id" binding:"required"`
}
