package schema

type FileDownloadRequest struct {
	//FileId uint64 `uri:"id" binding:"required"`
	Hash string `uri:"hash" binding:"required"`
}

//type FileUserDownloadRequest struct {
//	UserFileId uint64 `uri:"id" binding:"required"`
//}

//type UserIDTokenRequest struct {
//	IDToken string `form:"id_token"`
//}
