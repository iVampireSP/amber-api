package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

// AddChatImage godoc
// @Summary      添加图片
// @Description  将一个图片添加到聊天记录中
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Param        schema.ChatDownloadRemoteFileRequest  body  schema.ChatDownloadRemoteFileRequest false  "远程文件"
// @Param        image  formData  file  false  "图片"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/images [post]
func (u *ChatController) AddChatImage(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	// 检查状态是否是回复中
	isStreaming, err := u.isStreaming(c, chatRequest.ChatId)
	if err != nil || isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	var chatDownloadRemoteFileRequest = &schema.ChatDownloadRemoteFileRequest{}

	var uploaded = true
	var uploadedFile *multipart.FileHeader

	// 检查 formData 是否有 image，如果没有则尝试绑定结构体
	if c.ContentType() == "multipart/form-data" {
		var request = &schema.ChatMessageAddImageRequest{}
		err = c.ShouldBind(request)
		if err != nil {
			response.Status(http.StatusBadRequest).Error(err).Send()
			return
		}

		uploadedFile = request.Image
	} else {
		uploaded = false

		// 尝试绑定结构体
		err = c.ShouldBindJSON(chatDownloadRemoteFileRequest)
		if err != nil {
			//response.Status(http.StatusBadRequest).Error(err).Send()
			response.Status(http.StatusBadRequest).Error(consts.ErrFileUrlRequired).Send()
			return
		}
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}

	chatEntity, err := u.chatService.GetChat(c, chatRequest.ChatId)
	if err != nil || chatEntity.UserId != u.authService.GetUserId(c) {
		if errors.Is(err, consts.ErrChatNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	var file = &entity.File{}
	if uploaded {
		f, err := uploadedFile.Open()
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(consts.ErrUnableOpenFile).Send()
			return
		}

		file, err = u.fileService.CreateFile(c, f)
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
		file, err = u.fileService.CreateFileFromUrl(c, chatDownloadRemoteFileRequest.Url)
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

	if lastChatMessage.Role == schema.RoleFile && lastChatMessage.Content == file.Id.String() {
		response.Message(consts.HintProvideSameImage)
	}

	var chatMessage entity.ChatMessage
	chatMessage.ChatId = chatEntity.Id
	chatMessage.Content = file.Id.String()
	chatMessage.Role = schema.RoleFile

	err = u.cm.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(consts.ErrCreateChatMessage).Send()
		return
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
// @Param        schema.GetPublicChatMessageRequest  formData  schema.GetPublicChatMessageRequest false  "GetPublicChatMessageRequest"
// @Param        schema.ChatDownloadRemoteFileRequest  body  schema.ChatDownloadRemoteFileRequest false  "远程文件"
// @Param        image  formData  file  true  "图片"
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
	assistantShare, err := u.assistantService.GetShareByToken(c, getPublicChatMessageRequest.AssistantToken)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, getPublicChatMessageRequestParams.ChatId)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	// 检查 assistant id 是否一致
	if chatEntity.AssistantId != assistantShare.AssistantId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	if chatEntity.Owner != schema.OwnerGuest || (chatEntity.GuestId != nil && *chatEntity.GuestId != getPublicChatMessageRequest.GuestId) {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	// 检查状态是否是回复中
	isStreaming, err := u.isStreaming(c, chatEntity.Id)
	if err != nil || isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}
	var chatDownloadRemoteFileRequest = &schema.ChatDownloadRemoteFileRequest{}

	var uploaded = true
	var uploadedFile *multipart.FileHeader
	var file *entity.File

	// 检查 formData 是否有 image，如果没有则尝试绑定结构体
	if c.ContentType() == "multipart/form-data" {
		var request = &schema.ChatMessageAddImageRequest{}
		err = c.ShouldBind(request)
		if err != nil {
			response.Status(http.StatusBadRequest).Error(err).Send()
			return
		}

		uploadedFile = request.Image
	} else {
		uploaded = false

		// 尝试绑定结构体
		err = c.ShouldBindJSON(chatDownloadRemoteFileRequest)
		if err != nil {
			//response.Status(http.StatusBadRequest).Error(err).Send()
			response.Status(http.StatusBadRequest).Error(consts.ErrFileUrlRequired).Send()
			return
		}
	}

	if uploaded {
		f, err := uploadedFile.Open()
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(consts.ErrUnableOpenFile).Send()
			return
		}

		file, err = u.fileService.CreateFile(c, f)
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
		file, err = u.fileService.CreateFileFromUrl(c, chatDownloadRemoteFileRequest.Url)
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

	if lastChatMessage.Role == schema.RoleFile && lastChatMessage.Content == file.Id.String() {
		response.Message(consts.HintProvideSameImage)
	}

	var chatMessage entity.ChatMessage
	chatMessage.ChatId = chatEntity.Id
	chatMessage.Content = file.Id.String()
	chatMessage.Role = schema.RoleFile

	err = u.cm.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(consts.ErrCreateChatMessage).Send()
		return
	}

	response.Status(http.StatusOK).Data(chatMessageResponse).Send()
}
