package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"rag-new/internal/entity"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/chat"
	"rag-new/pkg/consts"
	"rag-new/pkg/random"
	"strconv"
)

type ChatController struct {
	authService *auth.Service
	chatService *chat.Service
	redis       *redis.Client
}

func NewChatController(authService *auth.Service, chatService *chat.Service, redis *redis.Client) *ChatController {
	return &ChatController{authService, chatService, redis}
}

// List godoc
// @Summary      获取所有 Chat
// @Description  get string by ID
// @Tags         chat
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
// @Tags         chat
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
// @Tags         chat
// @Accept       json
// @Produce      json
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

// ListChatMessage godoc
// @Summary      查看聊天记录
// @Description  get string by ID
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Success      200  {object}  schema.ResponseBody{data=[]entity.ChatMessage}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/messages [get]
func (u *ChatController) ListChatMessage(c *gin.Context) {
	var response = schema.NewResponse(c)

	chatId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, int64(chatId))
	if err != nil {
		if errors.Is(err, consts.ErrChatNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	chatHistories, err := u.chatService.GetChatMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(chatHistories).Send()
}

// AddChatMessage godoc
// @Summary      添加聊天记录
// @Description  get string by ID
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Chat ID"
// @Param        message  body  schema.ChatMessageAddRequest  true  "Message"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/messages [post]
func (u *ChatController) AddChatMessage(c *gin.Context) {
	var response = schema.NewResponse(c)

	chatIdStr := c.Param("id")
	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var request schema.ChatMessageAddRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}

	// 检测 chat 是否存在缓存
	cmd := u.redis.Get(c, u.getCacheKey("entity:"+chatIdStr))
	result, err := cmd.Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			response.Status(http.StatusInternalServerError).Error(cmd.Err()).Send()

			return
		}
	} else {
		chatMessageResponse.StreamId = result

		response.Status(http.StatusConflict).Error(consts.ErrChatStreamNotOpen).Data(chatMessageResponse).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, int64(chatId))
	if err != nil || chatEntity.UserId != u.authService.GetUserId(c) {
		if errors.Is(err, consts.ErrChatNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	// last chat message
	lastChatMessage, err := u.chatService.GetLatestMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if lastChatMessage.Role == entity.RoleHuman {
		lastChatMessage.Content = request.Message
		err := u.chatService.UpdateMessageContent(c, lastChatMessage)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		// 如果 stream id 过期了，但 role 还是 entity.RoleHuman ，则说明没有打开 chat stream，重新生成一个 stream id
		randomStreamId, err := u.generateChatStream(c, chatIdStr)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
		chatMessageResponse.StreamId = randomStreamId

		response.Status(http.StatusConflict).Error(consts.ErrChatStreamNotOpenAndOverrideMessage).Data(chatMessageResponse).Send()
		return
	}

	var chatMessage entity.ChatMessage
	chatMessage.ChatId = chatEntity.ID
	chatMessage.Content = request.Message
	chatMessage.Role = entity.RoleHuman

	err = u.chatService.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	randomStreamId, err := u.generateChatStream(c, chatIdStr)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	chatMessageResponse.StreamId = randomStreamId

	response.Status(http.StatusOK).Data(chatMessageResponse).Send()
}

func (u *ChatController) getCacheKey(key string) string {
	return fmt.Sprintf("chat:%s", key)
}

func (u *ChatController) generateChatStream(c context.Context, chatId string) (streamId string, err error) {
	var randomId = random.String(32)
	// 保存 chat stream id
	err = u.redis.Set(c, u.getCacheKey("entity:"+chatId), randomId, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	err = u.redis.Set(c, u.getCacheKey("stream:"+randomId), chatId, 0).Err()
	if err != nil {
		return "", err
	}

	return randomId, nil
}
