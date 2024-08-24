package schema

type FileDownloadRequest struct {
	FileId uint64 `uri:"id" binding:"required"`
}
