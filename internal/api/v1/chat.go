package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/base/redis"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/chat_message"
	"rag-new/internal/service/file"
	"rag-new/internal/service/library"
	"rag-new/internal/service/llm"
	"rag-new/internal/service/memory"
	"rag-new/internal/service/message_block"
	"rag-new/internal/service/tool"
	"rag-new/pkg/consts"
	"strconv"
)

const eventName = "data"
const eventDone = "[DONE]"

type ChatController struct {
	authService      *auth.Service
	chatService      *chat.Service
	redis            *redis.Redis
	llmService       *llm.Service
	logger           *logger.Logger
	assistantService *assistant.Service
	cm               *chat_message.Service
	fileService      *file.Service
	memoryService    *memory.Service
	libraryService   *library.Service
	config           *conf.Config
	toolService      *tool.Service
	messageBlock     *message_block.Service
}

func NewChatController(
	authService *auth.Service,
	chatService *chat.Service,
	redis *redis.Redis,
	llmService *llm.Service,
	logger *logger.Logger,
	assistantService *assistant.Service,
	chatMessageService *chat_message.Service,
	config *conf.Config,
	fileService *file.Service,
	memoryService *memory.Service,
	libraryService *library.Service,
	toolService *tool.Service,
	messageBlock *message_block.Service,
) *ChatController {
	return &ChatController{
		authService,
		chatService,
		redis,
		llmService,
		logger,
		assistantService,
		chatMessageService,
		fileService,
		memoryService,
		libraryService,
		config,
		toolService,
		messageBlock,
	}
}

// List godoc
// @Summary      获取所有 Chat
// @Description  列出当前账户下的所有的对话
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

		assistantEntity, err := u.assistantService.GetAssistant(c, schema.EntityId(assistantIdInt))
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
// @Description  创建一个对话，如果不指定 Assistant ID，将会使用默认 Assistant。默认 Assistant 不支持上传文件以及使用外部工具。
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

// Show godoc
// @Summary      显示一个对话的数据
// @Description  将返回一个实体
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.Chat}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id} [get]
func (u *ChatController) Show(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, chatRequest.ChatId)
	if err != nil {
		if errors.Is(err, consts.ErrChatNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	if chatEntity.Id == consts.NoRecord || chatEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrChatNotFound).Send()
		return
	}

	response.Status(http.StatusOK).Data(chatEntity).Send()
}

// Update godoc
// @Summary      更新对话
// @Description  可以重新设置对话的一些信息
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Param        schema.ChatUpdateRequest  body  schema.ChatUpdateRequest true  "ChatUpdateRequest"
// @Success      200  {object}  schema.ResponseBody{data=entity.Chat}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id} [put]
func (u *ChatController) Update(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var chatUpdateRequest = &schema.ChatUpdateRequest{}
	err = c.ShouldBindJSON(chatUpdateRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}
	chatEntity, err := u.chatService.GetChat(c, chatRequest.ChatId)
	if err != nil {
		if errors.Is(err, consts.ErrChatNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	if chatEntity.Id == consts.NoRecord || chatEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrChatNotFound).Send()
		return
	}

	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, chatEntity.Id)
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	chatEntity.Name = chatUpdateRequest.Name
	if chatUpdateRequest.ExpiredAt != nil {
		chatEntity.ExpiredAt = &chatUpdateRequest.ExpiredAt.Time
	} else {
		chatEntity.ExpiredAt = nil
	}

	chatEntity.Prompt = chatUpdateRequest.Prompt

	if chatUpdateRequest.AssistantId != nil {
		canUse, err := u.assistantService.CanUse(c, u.authService.GetUserId(c), *chatUpdateRequest.AssistantId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		if !canUse {
			response.Status(http.StatusForbidden).Error(consts.ErrAssistantNotPublic).Send()
			return
		}

		chatEntity.AssistantId = chatUpdateRequest.AssistantId

	} else {
		chatEntity.AssistantId = nil
	}

	err = u.chatService.UpdateChat(c, chatEntity)
	if err != nil {
		if errors.Is(err, consts.ErrChatNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	response.Status(http.StatusOK).Data(chatEntity).Send()
}

// Delete godoc
// @Summary      Delete Chat
// @Description  删除一个对话以及聊天记录
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
	isStreaming := u.isStreaming(c, schema.EntityId(chatId))
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	err = u.chatService.DeleteChatFromUserId(c, schema.EntityId(chatId), u.authService.GetUserId(c))
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

func (u *ChatController) getChatIdStreamingKey(chatId schema.EntityId) string {
	return u.getCacheKey("entity:" + strconv.Itoa(int(chatId)) + ":streaming")
}

func (u *ChatController) isStreaming(ctx context.Context, chatId schema.EntityId) bool {
	// 检查状态是否是回复中
	chatIdStreamingKey := u.getChatIdStreamingKey(chatId)
	i, err := u.redis.Client.Exists(ctx, chatIdStreamingKey).Result()
	if err != nil {
		return false
	}
	if i != consts.NoRecord {
		return true
	}

	return false
}
