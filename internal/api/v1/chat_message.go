package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"rag-new/pkg/random"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ListChatMessage godoc
// @Summary      查看聊天记录
// @Description  获取一个对话的所有聊天记录
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.ChatMessageList}
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

	var chatMessageListResponse []*entity.ChatMessageList

	chatHistories, err := u.cm.GetChatMessage(c, chatEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		}
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	for _, chatMessage := range chatHistories {
		var cmr = &entity.ChatMessageList{}
		cmr.Id = chatMessage.Id
		cmr.CreatedAt = chatMessage.CreatedAt
		cmr.UpdatedAt = chatMessage.UpdatedAt
		cmr.ChatId = chatMessage.ChatId
		cmr.AssistantId = chatMessage.AssistantId
		cmr.Role = chatMessage.Role
		cmr.Content = chatMessage.Content
		cmr.FileId = chatMessage.FileId
		cmr.UserFileId = chatMessage.UserFileId
		if chatMessage.File != nil {
			cmr.File = chatMessage.File
		}
		if chatMessage.UserFile != nil {
			cmr.UserFile = chatMessage.UserFile
		}

		cmr.Hidden = chatMessage.Hidden
		cmr.PromptTokens = chatMessage.PromptTokens
		cmr.CompletionTokens = chatMessage.CompletionTokens
		cmr.TotalTokens = chatMessage.TotalTokens

		if chatMessage.Assistant != nil {
			cmr.Assistant = &struct {
				Id   schema.EntityId `json:"id"`
				Name string          `json:"name"`
			}{
				Id:   chatMessage.Assistant.Id,
				Name: chatMessage.Assistant.Name,
			}
		}

		chatMessageListResponse = append(chatMessageListResponse, cmr)
	}

	response.Status(http.StatusOK).Data(chatMessageListResponse).Send()
}

// AddChatMessage godoc
// @Summary      添加聊天记录
// @Description  添加一条消息
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
	var request schema.ChatMessageAddRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}

	var userInfo = u.authService.GetUser(c)
	var publicUser = &schema.UserPublicInfo{
		Name:      userInfo.Token.Name,
		Id:        userInfo.Token.Sub,
		ChatOwner: schema.OwnerUser,
	}

	var assistantEntity *entity.Assistant
	if request.AssistantId != nil {
		// 检测 Assistant 是否属于用户
		assistantEntity, err = u.assistantService.GetAssistant(c, *request.AssistantId)
		if err != nil {
			if errors.Is(err, consts.ErrAssistantNotFound) {
				response.Status(http.StatusNotFound).Error(err).Send()
			} else {
				response.Status(http.StatusInternalServerError).Error(err).Send()
			}
			return
		}
		if assistantEntity.UserId != userInfo.Token.Sub {
			response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
			return
		}

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

	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, chatRequest.ChatId)
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	// 检测 chat 是否存在缓存
	cmd := u.redis.Get(c, u.getCacheKey("entity:"+chatEntity.Id.String()))
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
			randomStreamId, err := u.generateChatStream(c, chatEntity.Id, publicUser, request.Variables)
			if err != nil {
				response.Status(http.StatusInternalServerError).Error(err).Send()
				return
			}
			chatMessageResponse.StreamId = randomStreamId

			response.Status(http.StatusConflict).Error(consts.ErrChatStreamNotOpenAndOverrideMessage).Data(chatMessageResponse).Send()
			return
		} else if lastChatMessage.Role == schema.RoleAssistant {
			// 如果是 Assistant 消息，则开始采样记忆
			if request.Role == schema.RoleHuman {
				// 如果对话没有 Assistant，则默认启用记忆
				if chatEntity.Assistant == nil {
					u.addMemory(c, userInfo, request)
				} else if !chatEntity.Assistant.DisableDefaultPrompt {
					u.addMemory(c, userInfo, request)
				}
			}
		}
	}

	// 消息写入列表
	var chatMessages []entity.ChatMessage

	// 如果 Role 是 File
	if request.Role == schema.RoleFile {
		// 不需要串流，也不需要获取知识库
		needStream = false
	}

	var assistantId = chatEntity.AssistantId

	if request.AssistantId != nil {
		// 如果这个 assistant 不是用户的
		assistantEntity2, err := u.assistantService.GetAssistant(c, *request.AssistantId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		if assistantEntity2.UserId != userInfo.Token.Sub {
			response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
			return
		}

		assistantId = request.AssistantId
	}

	chatMessages = append(chatMessages, entity.ChatMessage{
		ChatId:      chatEntity.Id,
		AssistantId: assistantId,
		Content:     request.Message,
		Role:        request.Role,
	})

	// 检测是否存在知识库
	if needStream && assistantEntity != nil && assistantEntity.LibraryId != nil {
		libraryEntity, err := u.libraryService.GetLibrary(c, *assistantEntity.LibraryId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		// 从知识库获取内容，并添加到历史上下文
		libraryResults, err := u.libraryService.SearchLibrary(c, request.Message, libraryEntity)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		var chunkContent = ""
		// 将 libraryResults 拼接起来
		for _, libraryResult := range libraryResults {
			chunkContent += libraryResult.Content
		}

		// 添加知识库消息
		chatMessages = append(chatMessages, entity.ChatMessage{
			ChatId:      chatEntity.Id,
			AssistantId: &assistantEntity.Id,
			Content:     chunkContent,
			Role:        schema.RoleSystem,
		})
	}

	// TODO: 如果 request.Message 的大小超过了 1mb, 则转换为文件。转换之前应该先判断助理是否存在知识库
	// Update: 其实我也不知道这个要不要做,感觉做了意义也不大
	// 如果存在知识库，则将文件放入知识库中，否则将不处理
	//if len(request.Message) > 1024*1024 {
	//	// 转换为 file
	//
	//}

	for _, chatMessage := range chatMessages {
		err = u.cm.CreateChatMessage(c, &chatMessage)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	chatMessageResponse.Stream = needStream
	if needStream {
		randomStreamId, err := u.generateChatStream(c, chatEntity.Id, publicUser, request.Variables)
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

type ChatStreamCache struct {
	ChatId    schema.EntityId
	Variables map[string]string
}

func (u *ChatController) generateChatStream(c context.Context,
	chatId schema.EntityId,
	userPublic *schema.UserPublicInfo,
	variables map[string]string) (streamId string, err error) {
	var randomId = random.String(32)
	// 保存 chat stream id
	err = u.redis.Set(c, u.getCacheKey("entity:"+chatId.String()), randomId, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	var csc = ChatStreamCache{
		ChatId:    chatId,
		Variables: variables,
	}

	chatJson, err := sonic.MarshalString(csc)
	if err != nil {
		return "", err
	}

	err = u.redis.Set(c, u.getCacheKey("stream:"+randomId), chatJson, consts.ChatStreamExpire).Err()
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

func (u *ChatController) addMemory(c context.Context, userInfo *schema.User, request schema.ChatMessageAddRequest) {
	go func() {
		u.logger.Sugar.Info("memory service adding: ", request.Message)

		err := u.memoryService.Add(c, request.Message, userInfo.Token.Sub)
		if err != nil {
			u.logger.Sugar.Error("memory service add error: ", err)
		}
	}()
}
