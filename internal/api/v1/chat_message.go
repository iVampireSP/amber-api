package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"rag-new/pkg/random"
	"strconv"
)

// ListChatMessage godoc
// @Summary      查看聊天记录
// @Description  get string by ID
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.ChatMessage}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/messages [get]
func (u *ChatController) ListChatMessage(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, chatRequest.ChatId)
	if err != nil {
		if errors.Is(err, consts.ErrChatNotFound) {
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
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Param        message  body  schema.ChatMessageAddRequest  true  "Message"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/messages [post]
func (u *ChatController) AddChatMessage(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}
	var chatIdStr = fmt.Sprintf("%d", chatRequest.ChatId)
	var request schema.ChatMessageAddRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}

	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, chatRequest.ChatId)
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

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

	var needStream = true

	// 不允许添加文件
	if request.Role == schema.RoleFile {
		response.Status(http.StatusBadRequest).Error(consts.ErrRoleCanNotBeFile).Send()
		return
	}

	// 如果不是 human 或者 hide_human，则不需要回复
	if request.Role != schema.RoleHuman && request.Role != schema.RoleHideHuman {
		// 不需要生成 ID,直接添加
		needStream = false
	}

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
	var userIdStr = strconv.Itoa(int(u.authService.GetUserId(c)))

	var userInfo = u.authService.GetUser(c)
	var publicUser = &schema.UserPublicInfo{
		Name:      userInfo.Token.Name,
		Id:        userIdStr,
		ChatOwner: schema.OwnerUser,
	}

	// last chat message
	lastChatMessage, err := u.cm.GetLatestMessage(c, chatEntity)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if lastChatMessage != nil {
		// 如果有悬垂工具调用（要调用 tool，但是没有找到 tool response 的场景）
		if lastChatMessage.Role == schema.RoleToolCall {
			// 一般这种情况，肯定是工具调用失败了，或者是程序错误，所以这里补一个 tool role, 表明工具失败
			// 那么删掉最后一条消息即可
			err = u.cm.DeleteChatMessage(c, lastChatMessage)
			if err != nil {
				response.Status(http.StatusInternalServerError).Error(err).Send()
				return
			}
		} else if lastChatMessage.Role == schema.RoleHuman {
			// 如果上一条消息是 Human 消息，则说明消息没有成功发送，覆盖上一条消息
			lastChatMessage.Content = request.Message
			err = u.cm.UpdateMessageContent(c, lastChatMessage)
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
	}

	var chatMessage entity.ChatMessage
	chatMessage.ChatId = chatEntity.Id
	chatMessage.Content = request.Message
	// 准备检测 Role
	chatMessage.Role = request.Role

	// 如果 Role 是 File
	if request.Role == schema.RoleFile {
		needStream = false
	}

	err = u.cm.CreateChatMessage(c, &chatMessage)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	chatMessageResponse.Stream = needStream
	if needStream {
		randomStreamId, err := u.generateChatStream(c, chatIdStr, publicUser)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
		chatMessageResponse.StreamId = randomStreamId
	}

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

// ClearChatMessage godoc
// @Summary      清空聊天记录
// @Description  清空当前聊天记录
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Success      200
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/clear [post]
func (u *ChatController) ClearChatMessage(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, chatRequest.ChatId)
	if err != nil {
		if errors.Is(err, consts.ErrChatNotFound) {
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
