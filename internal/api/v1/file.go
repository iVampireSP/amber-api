package v1

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"rag-new/internal/base/logger"
	"rag-new/internal/schema"
	"rag-new/internal/service/file"
	"rag-new/pkg/consts"
	"strconv"
)

type FileController struct {
	fileService *file.Service
	logger      *logger.Logger
}

func NewFileController(fileService *file.Service, logger *logger.Logger) *FileController {
	return &FileController{fileService, logger}
}

// DownloadFile godoc
// @Summary      下载文件
// @Description  根据 File ID 下载文件
// @Tags         file
// @Accept       json
// @Produce      json
// @Param        schema.FileDownloadRequest  path  schema.FileDownloadRequest true  "File ID"
// @Router       /api/v1/files/{id}/download [get]
func (f *FileController) DownloadFile(c *gin.Context) {
	var response = schema.NewResponse(c)

	var fileDownloadRequest = &schema.FileDownloadRequest{}
	if err := c.ShouldBindUri(fileDownloadRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	fileExists, err := f.fileService.ExistsFileById(c, schema.EntityId(fileDownloadRequest.FileId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if !fileExists {
		response.Status(http.StatusNotFound).Error(consts.ErrFileNotExists).Send()
		return
	}

	fileEntity, err := f.fileService.GetFileById(c, schema.EntityId(fileDownloadRequest.FileId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
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
