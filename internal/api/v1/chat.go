package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/chat_message"
	"rag-new/internal/service/file"
	"rag-new/internal/service/llm"
	"rag-new/pkg/consts"
	"strconv"
)

const eventName = "data"
const eventDone = "[DONE]"

type ChatController struct {
	authService      *auth.Service
	chatService      *chat.Service
	redis            *redis.Client
	llmService       *llm.Service
	logger           *logger.Logger
	assistantService *assistant.Service
	cm               *chat_message.Service
	fileService      *file.Service
	config           *conf.Config
}

func NewChatController(authService *auth.Service,
	chatService *chat.Service, redis *redis.Client, llmService *llm.Service, logger *logger.Logger, assistantService *assistant.Service, chatMessageService *chat_message.Service, config *conf.Config, fileService *file.Service) *ChatController {
	return &ChatController{authService, chatService, redis, llmService, logger, assistantService, chatMessageService, fileService, config}
}

// List godoc
// @Summary      获取所有 Chat
// @Description  get string by ID
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        assistant_id  query  int  false  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Chat}
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/chats [get]
func (u *ChatController) List(c *gin.Context) {
	var response = schema.NewResponse(c)

	// 检测 query 中是否有 assistant_id
	assistantId, _ := c.GetQuery("assistant_id")
	if assistantId != "" {
		assistantIdInt, err := strconv.Atoi(assistantId)
		if err != nil {
			response.Status(http.StatusBadRequest).Error(err).Send()
			return
		}

		assistantEntity, err := u.assistantService.GetAssistant(c, int64(assistantIdInt))
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
		if assistantEntity.Id == consts.NoRecord || assistantEntity.UserId != u.authService.GetUserId(c) {
			response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
			return
		}

		chatEntities, err := u.chatService.ListChatFromAssistantIdWithAssistant(c, assistantEntity)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
		response.Status(http.StatusOK).Data(chatEntities).Send()
	} else {
		chatEntities, err := u.chatService.ListChatFromUserId(c, u.authService.GetUserId(c))
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		response.Status(http.StatusOK).Data(chatEntities).Send()
	}

}

// Create godoc
// @Summary      Create Chat
// @Description  get string by ID
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        chat  body  	schema.ChatCreateRequest  true  "Chat"
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
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Chat ID"
// @Success      200  {object}  schema.ResponseBody
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

	// 检查状态是否是回复中
	isStreaming, err := u.isStreaming(c, int64(chatId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
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

func (u *ChatController) getChatIdStreamingKey(chatId int64) string {
	return u.getCacheKey("entity:" + strconv.Itoa(int(chatId)) + ":streaming")
}

func (u *ChatController) isStreaming(ctx context.Context, chatId int64) (bool, error) {
	// 检查状态是否是回复中
	chatIdStreamingKey := u.getChatIdStreamingKey(chatId)
	i, err := u.redis.Exists(ctx, chatIdStreamingKey).Result()
	if err != nil {
		return false, err
	}
	if i != consts.NoRecord {
		return true, nil
	}

	return false, nil
}
