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
// @Param        image  formData  file  true  "图片"
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

	var request schema.ChatMessageAddImageRequest
	err = c.ShouldBind(&request)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
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

	f, err := request.Image.Open()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(consts.ErrUnableOpenFile).Send()
		return
	}

	file, err := u.fileService.CreateFile(c, f)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var chatMessage entity.ChatMessage
	chatMessage.ChatId = chatEntity.Id
	chatMessage.Content = file.Id.String()
	chatMessage.Role = schema.RoleImage

	err = u.cm.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(consts.ErrCreateChatMessage).Send()
		return
	}

	response.Status(http.StatusOK).Data(chatMessageResponse).Send()

	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			u.logger.Sugar.Error(err)
			return
		}
	}(f)
}
