package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/chat"
	"rag-new/pkg/consts"
	"strconv"
)

type ChatController struct {
	authService *auth.Service
	chatService *chat.Service
}

func NewChatController(authService *auth.Service, chatService *chat.Service) *ChatController {
	return &ChatController{authService, chatService}
}

// List godoc
// @Summary      获取所有 Chat
// @Description  get string by ID
// @Tags         ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Chat}
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/chats [get]
func (u *ChatController) List(c *gin.Context) {
	var response = schema.NewResponse(c)
	chatEntities, err := u.chatService.ListChatFromUserId(c, u.authService.GetUserId(c))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(chatEntities).Send()
}

// Create godoc
// @Summary      Create Chat
// @Description  get string by ID
// @Tags         ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  schema.ResponseBody{data=entity.Chat}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats [post]
func (u *ChatController) Create(c *gin.Context) {
	var response = schema.NewResponse(c)
	var createChatReq = schema.ChatCreateRequest{}
	if err := c.ShouldBindJSON(&createChatReq); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	createChatReq.UserId = u.authService.GetUserId(c)

	chatEntity, err := u.chatService.CreateChat(c, &createChatReq)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(chatEntity).Send()
}

// Delete godoc
// @Summary      Delete Chat
// @Description  get string by ID
// @Tags         ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  schema.ResponseBody{data=schema.CurrentUserResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id} [delete]
func (u *ChatController) Delete(c *gin.Context) {
	var response = schema.NewResponse(c)

	chatId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	err = u.chatService.DeleteChatFromUserId(c, int64(chatId), u.authService.GetUserId(c))
	if err != nil {
		if errors.Is(err, consts.ErrChatNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	response.Status(http.StatusOK).Send()
}
