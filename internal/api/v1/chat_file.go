package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

// AddChatFile godoc
// @Summary      添加文件
// @Description  将一个文件添加到聊天记录中
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Param        schema.ChatDownloadRemoteFileRequest  formData  schema.ChatDownloadRemoteFileRequest false  "远程文件"
// @Param        file  formData  file  false  "文件"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/files [post]
func (u *ChatController) AddChatFile(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, chatRequest.ChatId)
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	var chatDownloadRemoteFileRequest = &schema.ChatDownloadRemoteFileRequest{}

	var uploaded = true
	var uploadedFile *multipart.FileHeader

	// 检查 formData 是否有 文件，如果没有则尝试绑定结构体
	if c.ContentType() == "multipart/form-data" {
		var request = &schema.ChatMessageAddFileRequest{}
		err = c.ShouldBind(request)
		if err != nil {
			response.Status(http.StatusBadRequest).Error(err).Send()
			return
		}

		uploadedFile = request.File
	} else {
		uploaded = false

		// 尝试绑定结构体
		err = c.ShouldBind(chatDownloadRemoteFileRequest)
		if err != nil {
			//response.Status(http.StatusBadRequest).Error(err).Send()
			response.Status(http.StatusBadRequest).Error(consts.ErrFileUrlRequired).Send()
			return
		}
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}
	var userId = u.authService.GetUserId(c)

	chatEntity, err := u.chatService.GetChat(c, chatRequest.ChatId)
	if err != nil || chatEntity.UserId != userId {
		if errors.Is(err, consts.ErrChatNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	var filename string
	var file = &entity.File{}
	if uploaded {
		if uploadedFile == nil {
			response.Status(http.StatusBadRequest).Error(consts.ErrFileRequired).Send()
			return
		}

		// 获取 filename
		filename = uploadedFile.Filename

		f, err := uploadedFile.Open()
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(consts.ErrUnableOpenFile).Send()
			return
		}

		file, err = u.fileService.CreateFile(c, f, false)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		defer func(f multipart.File) {
			err := f.Close()
			if err != nil {
				u.logger.Sugar.Error(err)
				return
			}
		}(f)
	} else {
		filename = "RemoteFile_" + chatDownloadRemoteFileRequest.Url
		file, err = u.fileService.CreateFileFromUrl(c, chatDownloadRemoteFileRequest.Url, false)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	// 如果 filename 过长(>254)，则截断
	if len(filename) > 254 {
		filename = filename[:254]
	}

	// last chat message
	lastChatMessage, err := u.cm.GetLatestMessage(c, chatEntity)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if lastChatMessage != nil {
		// 检测 FileId 或 UserFileId
		if lastChatMessage.Role == schema.RoleFile {
			if lastChatMessage.FileId != nil {
				if *lastChatMessage.FileId == file.Id {
					response.Status(http.StatusConflict).Error(consts.ErrProvideSameImage).Message(consts.ErrProvideSameImage.Error()).Send()
					return
				}
			}

			if lastChatMessage.UserFileId != nil {
				if lastChatMessage.UserFile.FileId == file.Id {
					response.Status(http.StatusConflict).Error(consts.ErrProvideSameImage).Message(consts.ErrProvideSameImage.Error()).Send()
					return
				}
			}
		}
	}

	// bind a file to user
	userFile, err := u.fileService.BindFileToUser(c, file, userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var chatMessage entity.ChatMessage
	chatMessage.ChatId = chatEntity.Id
	//chatMessage.Content = file.Id.String()
	chatMessage.Role = schema.RoleFile
	//chatMessage.FileId = &file.Id
	chatMessage.UserFileId = &userFile.Id
	chatMessage.UserFile = userFile

	err = u.cm.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(consts.ErrCreateChatMessage).Send()
		return
	}

	if u.libraryService.CanChunk(file) {
		// TODO: 如果设置了助理使用的知识库，则不适用默认知识库
		library, err := u.libraryService.DefaultLibrary(c, userId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		// 检测用户下有没有相同 File ID 的 Document，如果有，则比较文件是否一致
		document, err := u.libraryService.GetDocumentByFileAndLibrary(c, file, library)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		if document != nil {
			if *document.FileHash == file.FileHash {
				response.Status(http.StatusOK).Data(chatMessageResponse).Send()
				return
			}
		} else {
			document = &entity.Document{
				LibraryId: library.Id,
				FileId:    &file.Id,
				FileHash:  &file.FileHash,
				Name:      filename,
			}
		}

		err = u.libraryService.CreateDocument(c, document)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		// Chunk
		go func() {
			err := u.libraryService.ChunkFileToDocument(c, file, document)
			if err != nil {
				u.logger.Sugar.Error(err)
			}
		}()
	}

	response.Status(http.StatusOK).Data(chatMessageResponse).Send()

}

// AddPublicChatImage godoc
// @Summary      添加图片
// @Description  将一个图片添加到聊天记录中
// @Tags         chat_public
// @Accept       json
// @Produce      json
// @Param        schema.GetPublicChatMessageRequestParams  path  schema.GetPublicChatMessageRequestParams  true  "GetPublicChatMessageRequestParams"
// @Param        schema.GetPublicChatMessageRequest  formData  schema.GetPublicChatMessageRequest true  "GetPublicChatMessageRequest"
// @Param        schema.ChatDownloadRemoteFileRequest  formData  schema.ChatDownloadRemoteFileRequest false  "远程文件"
// @Param        image  formData  file  false  "图片"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chat_public/{chat_id}/images [post]
func (u *ChatController) AddPublicChatImage(c *gin.Context) {
	var response = schema.NewResponse(c)

	var getPublicChatMessageRequestParams schema.GetPublicChatMessageRequestParams
	if err := c.ShouldBindUri(&getPublicChatMessageRequestParams); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var getPublicChatMessageRequest schema.GetPublicChatMessageRequest
	if err := c.ShouldBind(&getPublicChatMessageRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	// get assistant by token
	next, chatEntity := u.canPublicChatNext(c, response, &getPublicChatMessageRequest, &getPublicChatMessageRequestParams)
	if !next {
		return
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}
	var chatDownloadRemoteFileRequest = schema.ChatDownloadRemoteFileRequest{}

	var uploaded = true
	var uploadedFile *multipart.FileHeader

	// 检查 formData 是否有 image，如果没有则尝试绑定结构体
	if c.ContentType() == "multipart/form-data" {
		var request = &schema.ChatMessageAddImageRequest{}
		err := c.ShouldBind(request)
		if err != nil {
			response.Status(http.StatusBadRequest).Error(err).Send()
			return
		}

		uploadedFile = request.Image
	} else {
		uploaded = false

		// 尝试绑定结构体
		err := c.ShouldBind(chatDownloadRemoteFileRequest)
		if err != nil {
			//response.Status(http.StatusBadRequest).Error(err).Send()
			response.Status(http.StatusBadRequest).Error(consts.ErrFileUrlRequired).Send()
			return
		}
	}

	var file *entity.File
	if uploaded {
		f, err := uploadedFile.Open()
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(consts.ErrUnableOpenFile).Send()
			return
		}

		file, err = u.fileService.CreateFile(c, f, true)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		defer func(f multipart.File) {
			err := f.Close()
			if err != nil {
				u.logger.Sugar.Error(err)
				return
			}
		}(f)
	} else {
		var err error
		file, err = u.fileService.CreateFileFromUrl(c, chatDownloadRemoteFileRequest.Url, true)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	// last chat message
	lastChatMessage, err := u.cm.GetLatestMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if lastChatMessage != nil {
		// 检测 FileId 或 UserFileId
		if lastChatMessage.Role == schema.RoleFile {
			if lastChatMessage.FileId != nil {
				if *lastChatMessage.FileId == file.Id {
					response.Status(http.StatusConflict).Error(consts.ErrProvideSameImage).Message(consts.ErrProvideSameImage.Error()).Send()
					return
				}
			}

			if lastChatMessage.UserFileId != nil {
				if lastChatMessage.UserFile.FileId == file.Id {
					response.Status(http.StatusConflict).Error(consts.ErrProvideSameImage).Message(consts.ErrProvideSameImage.Error()).Send()
					return
				}
			}
		}
	}
	//if lastChatMessage.Role == schema.RoleFile && *lastChatMessage.FileId == file.Id {
	//	response.Status(http.StatusConflict).Error(consts.ErrProvideSameImage).Message(consts.ErrProvideSameImage.Error()).Send()
	//	return
	//}

	var chatMessage entity.ChatMessage
	chatMessage.ChatId = chatEntity.Id
	//chatMessage.Content = file.Id.String()
	chatMessage.Role = schema.RoleFile
	chatMessage.FileId = &file.Id

	err = u.cm.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(consts.ErrCreateChatMessage).Send()
		return
	}

	response.Status(http.StatusOK).Data(chatMessageResponse).Send()
}

func (u *ChatController) canPublicChatNext(c *gin.Context, response *schema.HttpResponse, getPublicChatMessageRequest *schema.GetPublicChatMessageRequest, getPublicChatMessageRequestParams *schema.GetPublicChatMessageRequestParams) (bool, *entity.Chat) {
	assistantShare, err := u.assistantService.GetShareByToken(c, getPublicChatMessageRequest.AssistantToken)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return false, nil
	}

	chatEntity, err := u.chatService.GetChat(c, getPublicChatMessageRequestParams.ChatId)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return false, nil
	}

	// 检查 assistant id 是否一致
	if *chatEntity.AssistantId != assistantShare.AssistantId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return false, nil
	}

	if chatEntity.Owner != schema.OwnerGuest || (chatEntity.GuestId != nil && *chatEntity.GuestId != getPublicChatMessageRequest.GuestId) {
		response.Status(http.StatusForbidden).Error(err).Send()
		return false, nil
	}

	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, chatEntity.Id)
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return false, nil
	}

	return true, chatEntity
}
