package v1

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"rag-new/internal/base/logger"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/file"
	"rag-new/pkg/consts"
	"strconv"
	"strings"
)

type FileController struct {
	fileService *file.Service
	authService *auth.Service
	logger      *logger.Logger
}

func NewFileController(fileService *file.Service, logger *logger.Logger, authService *auth.Service) *FileController {
	return &FileController{fileService, authService, logger}
}

// DownloadImage godoc
// @Summary      下载图片
// @Description  根据文件 Hash 下载文件。仅支持图片下载，且图片具有有效期
// @Tags         file
// @Accept       json
// @Produce      json
// @Param        schema.FileDownloadRequest  path  schema.FileDownloadRequest true  "File ID"
// @Success 	 200 {file} file
// @Router       /api/v1/files/download/{hash} [get]
func (f *FileController) DownloadImage(c *gin.Context) {
	var response = schema.NewResponse(c)

	var fileDownloadRequest = &schema.FileDownloadRequest{}
	if err := c.ShouldBindUri(fileDownloadRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	fileEntity, err := f.fileService.GetFileByFileHash(c, fileDownloadRequest.Hash)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// 如果 mimetype 不为 image/ 则 403
	if fileEntity.MimeType == "" {
		response.Status(http.StatusForbidden).Error(consts.ErrFileNotImage).Send()
		return
	}

	// 检测是不是 image/ 开头
	if !strings.HasPrefix(fileEntity.MimeType, "image/") {
		response.Status(http.StatusForbidden).Error(consts.ErrFileNotImage).Send()
		return
	}

	size, bucketFile, err := f.fileService.GetBucketFile(c, fileEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// 设置 mime type
	//c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileEntity.FileHash)
	c.Writer.Header().Set("Content-Type", fileEntity.MimeType)
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	c.Writer.Header().Set("Cache-Control", "public, max-age=86400")
	c.Writer.Flush()

	defer func(bucketFile io.ReadCloser) {
		err := bucketFile.Close()
		if err != nil {
			f.logger.Sugar.Error(err)

			return
		}
	}(bucketFile)

	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, bucketFile)
	if err != nil {
		f.logger.Sugar.Error(err)

		return
	}

}

//
//// DownloadUserFile godoc
//// @Summary      下载用户开文件
//// @Description  根据 File ID 下载文件。如果文件是私有的，将无法下载
//// @Tags         file
//// @Accept       json
//// @Produce      json
//// @Param        schema.FileUserDownloadRequest  path  schema.FileUserDownloadRequest true  "User File ID"
//// @Param        schema.UserIDTokenRequest  query  schema.UserIDTokenRequest true  "ID Token"
//// @Success 	 200 {file} file
//// @Router       /api/v1/files/user/{id}/download [get]
//func (f *FileController) DownloadUserFile(c *gin.Context) {
//	var response = schema.NewResponse(c)
//
//	var fileUserDownloadRequest = &schema.FileUserDownloadRequest{}
//	if err := c.ShouldBindUri(fileUserDownloadRequest); err != nil {
//		response.Status(http.StatusBadRequest).Error(err).Send()
//		return
//	}
//
//	var userIDTokenRequest = &schema.UserIDTokenRequest{}
//	if err := c.ShouldBindQuery(userIDTokenRequest); err != nil {
//		response.Status(http.StatusBadRequest).Error(err).Send()
//		return
//	}
//
//	user, err := f.authService.GetUserFromIdToken(userIDTokenRequest.IDToken)
//	if err != nil {
//		response.Status(http.StatusInternalServerError).Error(err).Send()
//		return
//	}
//
//	userFileExists, err := f.fileService.ExistsUserFileById(c, schema.EntityId(fileUserDownloadRequest.UserFileId))
//	if err != nil {
//		response.Status(http.StatusInternalServerError).Error(err).Send()
//		return
//	}
//	if !userFileExists {
//		response.Status(http.StatusNotFound).Error(consts.ErrFileNotExists).Send()
//		return
//	}
//
//	userFileEntity, err := f.fileService.GetUserFile(c, schema.EntityId(fileUserDownloadRequest.UserFileId))
//	if err != nil {
//		response.Status(http.StatusInternalServerError).Error(err).Send()
//		return
//	}
//
//	if userFileEntity.UserId != user.Token.Sub {
//		response.Status(http.StatusNotFound).Error(consts.ErrFileNotExists).Send()
//		return
//	}
//
//	size, bucketFile, err := f.fileService.GetBucketFile(c, userFileEntity.File)
//	if err != nil {
//		response.Status(http.StatusInternalServerError).Error(err).Send()
//		return
//	}
//
//	c.Writer.Header().Set("Content-Type", userFileEntity.File.MimeType)
//	c.Writer.Header().Set("Content-Length", strconv.FormatInt(size, 10))
//	c.Writer.Flush()
//
//	defer func(bucketFile io.ReadCloser) {
//		err := bucketFile.Close()
//		if err != nil {
//			f.logger.Sugar.Error(err)
//
//			return
//		}
//	}(bucketFile)
//
//	c.Status(http.StatusOK)
//	_, err = io.Copy(c.Writer, bucketFile)
//	if err != nil {
//		f.logger.Sugar.Error(err)
//
//		return
//	}
//
//}
