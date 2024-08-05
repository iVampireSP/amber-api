package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"rag-new/internal/base/logger"
	"rag-new/internal/entity"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/chat_message"
	"rag-new/internal/service/llm"
	"rag-new/pkg/consts"
	"rag-new/pkg/random"
	"strconv"
)

const eventName = "data"

type ChatController struct {
	authService      *auth.Service
	chatService      *chat.Service
	redis            *redis.Client
	llmService       *llm.Service
	logger           *logger.Logger
	assistantService *assistant.Service
	cm               *chat_message.Service
}

func NewChatController(authService *auth.Service,
	chatService *chat.Service, redis *redis.Client, llmService *llm.Service, logger *logger.Logger, assistantService *assistant.Service, chatMessageService *chat_message.Service) *ChatController {
	return &ChatController{authService, chatService, redis, llmService, logger, assistantService, chatMessageService}
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
		if assistantEntity.ID == consts.NoRecord || assistantEntity.UserId != u.authService.GetUserId(c) {
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

// ListChatMessage godoc
// @Summary      查看聊天记录
// @Description  get string by ID
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Chat ID"
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

	if chatEntity.ID == consts.NoRecord || chatEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrChatNotFound).Send()
		return
	}

	chatHistories, err := u.cm.GetChatMessage(c, chatEntity)
	//chatHistories, err := u.cm.GetChatMessageWithHide(c, chatEntity)
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
// @Security     ApiKeyAuth
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
	lastChatMessage, err := u.cm.GetLatestMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var userIdStr = strconv.Itoa(int(u.authService.GetUserId(c)))

	var userInfo = u.authService.GetUser(c)
	var publicUser = &schema.UserPublicInfo{
		Name:      userInfo.Token.Name,
		Id:        userIdStr,
		ChatOwner: schema.OwnerUser,
	}

	if lastChatMessage.Role == schema.RoleHuman {
		lastChatMessage.Content = request.Message
		err := u.cm.UpdateMessageContent(c, lastChatMessage)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		// 如果 stream id 过期了，但 role 还是 entity.RoleHuman ，则说明没有打开 chat stream，重新生成一个 stream id
		randomStreamId, err := u.generateChatStream(c, chatIdStr, publicUser)
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
	chatMessage.Role = schema.RoleHuman

	err = u.cm.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	randomStreamId, err := u.generateChatStream(c, chatIdStr, publicUser)
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

func (u *ChatController) generateChatStream(c context.Context, chatId string, userPublic *schema.UserPublicInfo) (streamId string, err error) {
	var randomId = random.String(32)
	// 保存 chat stream id
	err = u.redis.Set(c, u.getCacheKey("entity:"+chatId), randomId, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	err = u.redis.Set(c, u.getCacheKey("stream:"+randomId), chatId, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	userJson, err := sonic.Marshal(userPublic)
	if err != nil {
		return "", err
	}

	err = u.redis.Set(c, u.getCacheKey("stream:"+randomId+":user"), userJson, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	return randomId, nil
}

// Stream godoc
// @Summary      流式传输聊天内容
// @Description  get string by ID
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     none
// @Param        id  path  int  true  "Chat ID"
// @Param        stream_id  path  string  true  "Chat stream id"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/stream/{stream_id} [get]
func (u *ChatController) Stream(c *gin.Context) {
	var response = schema.NewResponse(c)
	// 检查 stream id 是否存在
	streamIdStr := c.Param("stream_id")
	streamIdCacheKey := u.getCacheKey("stream:" + streamIdStr)
	// 检查缓存是否存在
	i, err := u.redis.Exists(c, streamIdCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if i == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatStreamNotFound).Send()
		return
	}

	// 获取 chat id
	chatIdStr, err := u.redis.Get(c, streamIdCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var chatIdStreamKey = u.getChatIdStreamingKey(int64(chatId))
	// 检查状态是否是回复中
	isStreaming, err := u.isStreaming(c, int64(chatId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if isStreaming {
		response.Status(http.StatusConflict).Error(consts.ErrChatStreamingPleaseWait).Send()
		return
	} else {
		// 不在回复中，则设置缓存键，防止再次请求
		cmd := u.redis.Set(c, chatIdStreamKey, 1, consts.ChatStreamExpire)
		if cmd.Err() != nil {
			response.Status(http.StatusInternalServerError).Error(cmd.Err()).Send()
			return
		}
	}

	chatEntity, err := u.chatService.GetChat(c, int64(chatId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if chatEntity.ID == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatNotFound).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, chatEntity.AssistantId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// 获取 assistant 绑定的 tools
	tools, err := u.assistantService.ToLLMTool(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// 提取 history
	histories, err := u.cm.GetChatMessageWithHide(c, chatEntity)
	var llmResponseChan = make(chan *llm.AssistantResponse)

	// SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	streamUserCacheKey := u.getCacheKey("stream:" + streamIdStr + ":user")

	// 检查缓存是否存在
	i, err = u.redis.Exists(c, streamUserCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if i == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatStreamNotFound).Send()
		return
	}

	userCmd, err := u.redis.Get(c, streamUserCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	user := &schema.UserPublicInfo{}
	err = sonic.Unmarshal([]byte(userCmd), user)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var prompt = `
Your name: ` + assistantEntity.Name + `current user give you` + `
Your description: ` + assistantEntity.Description + "(current user given)"
	if user != nil {
		prompt += `
current user's name: ` + user.Name + `(system hint you this)` + `
current user's id: ` + user.Id + "(system hint you this, user can't change it)"

	}

	if assistantEntity.Prompt != "" {
		prompt += "\n" + assistantEntity.Prompt
	}

	go func() {
		err = u.llmService.StreamChat(llmResponseChan, prompt, user, histories, tools...)
		if err != nil {
			u.logger.Sugar.Error(err)
			// 关闭连接
			c.AbortWithStatus(http.StatusInternalServerError)
			return
			//response.Status(http.StatusInternalServerError).Message("Streaming chat failed").Error(err).Send()
		}
	}()

	var llmFullMessage = ""
	var toolResponse = ""
	c.Stream(func(w io.Writer) bool {
		// Emit Server Sent Events compatible
		msg, ok := <-llmResponseChan
		if !ok {
			return false
		}

		if msg == nil {
			return true
		}

		j, err := sonic.Marshal(msg)
		if err != nil {
			u.logger.Sugar.Error(err)
		}

		c.SSEvent(eventName, string(j))

		c.Writer.Flush()

		switch msg.State {
		case llm.StateChunk:
			llmFullMessage += msg.ChunkMessage.Content
			return true
		case llm.StateToolSuccess:
			return true
		case llm.StateToolResponse:
			toolResponse += `Tool/Function Call Logs:
Function Name: ` + msg.ToolResponseMessage.FunctionName + `
Function Response: ` + msg.ToolResponseMessage.Content + `

`
			return true
		case llm.StateFailed:
			return false
		case llm.StateFinished:
			return false
		case llm.StateToolFailed:
			return false
		default:
			return true
		}
	})

	// close sse stream
	c.SSEvent("close", "")
	c.Writer.Flush()

	close(llmResponseChan)

	// 移除缓存
	u.redis.Del(c, streamIdCacheKey)
	u.redis.Del(c, u.getCacheKey("entity:"+chatIdStr))
	u.redis.Del(c, u.getCacheKey("stream:"+streamIdStr+":user"))
	u.redis.Del(c, chatIdStreamKey)

	var newMessage *entity.ChatMessage
	if toolResponse != "" {
		// 添加到消息 entity.ChatMessage
		newMessage = &entity.ChatMessage{
			Role:    schema.RoleHideSystem,
			Content: toolResponse,
			ChatId:  chatEntity.ID,
		}

		err = u.cm.CreateChatMessage(c, newMessage)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	if llmFullMessage != "" {
		// 添加到消息 entity.ChatMessage
		newMessage = &entity.ChatMessage{
			Role:    schema.RoleAssistant,
			Content: llmFullMessage,
			ChatId:  chatEntity.ID,
		}

		err = u.cm.CreateChatMessage(c, newMessage)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

	}

	response.Status(http.StatusOK).Send()
}

// ClearChatMessage godoc
// @Summary      清空聊天记录
// @Description  清空当前聊天记录
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Chat ID"
// @Success      200
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/clear [post]
func (u *ChatController) ClearChatMessage(c *gin.Context) {
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

	if chatEntity.ID == consts.NoRecord || chatEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrChatNotFound).Send()
		return
	}

	// 检查状态是否是回复中
	isStreaming, err := u.isStreaming(c, chatEntity.ID)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if isStreaming {
		response.Status(http.StatusConflict).Error(consts.ErrChatStreaming).Send()
		return
	}

	err = u.cm.ClearChatMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
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
