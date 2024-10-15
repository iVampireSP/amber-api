package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/iVampireSP/pkg/page"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"slices"
)

var allowedChatMessageRoles = []schema.ChatRole{
	schema.RoleHuman,
	schema.RoleHumanLater,
	schema.RoleHideHuman,
	schema.RoleSystem,
	schema.RoleHideSystem,
	schema.RoleSystemOverride,
	schema.RoleAssistant,
}

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
		//cmr.UserFileId = chatMessage.UserFileId
		if chatMessage.File != nil {
			cmr.File = chatMessage.File
		}
		//if chatMessage.UserFile != nil {
		//	cmr.UserFile = chatMessage.UserFile
		//}

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

// ListChatMessagePaginate godoc
// @Summary      分页获取聊天记录
// @Description  获取一个对话的所有聊天记录
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ChatRequest  path  schema.ChatRequest true  "Chat ID"
// @Param        schema.PaginationRequest  query  schema.PaginationRequest true  "schema.PaginationRequest"
// @Success      200  {object}  schema.ResponseBody{data=page.PagedResult[entity.ChatMessageList]}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chats/{id}/messages/paginate [get]
func (u *ChatController) ListChatMessagePaginate(c *gin.Context) {
	var response = schema.NewResponse(c)

	var chatRequest = &schema.ChatRequest{}
	err := c.ShouldBindUri(chatRequest)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var paginationRequest = &schema.PaginationRequest{}
	if err := c.ShouldBind(paginationRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var pagedResult = page.NewPagedResult[*entity.ChatMessageList]()

	pagedResult.Page = paginationRequest.Page
	pagedResult.PageSize = 2

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

	chatHistories, _, count, err := u.cm.GetChatMessagePageAsc(c, chatEntity, paginationRequest.Page, pagedResult.PageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		}
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	pagedResult.TotalCount = count

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
		//cmr.UserFileId = chatMessage.UserFileId
		if chatMessage.File != nil {
			cmr.File = chatMessage.File
		}
		//if chatMessage.UserFile != nil {
		//	cmr.UserFile = chatMessage.UserFile
		//}

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

		pagedResult.Add(cmr)
	}

	response.Status(http.StatusOK).Data(pagedResult.Output()).Send()
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

	if !slices.Contains(allowedChatMessageRoles, request.Role) {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatRoleNotAllowed).Send()
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

		// 检测是否公开
		if !assistantEntity.Public && assistantEntity.UserId != userInfo.Token.Sub {
			response.Status(http.StatusForbidden).Error(consts.ErrAssistantNotPublic).Send()
			return
		}

		// 检测是不是收藏的
		hasFavorite, err := u.assistantService.HasFavoriteAssistant(c, userInfo.Token.Sub, assistantEntity)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		if !hasFavorite {
			if assistantEntity.UserId != userInfo.Token.Sub {
				response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
				return
			}
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

	// 如果没有 Assistant，则获取对话的 Assistant
	if assistantEntity == nil && chatEntity.Assistant != nil {
		assistantEntity = chatEntity.Assistant
	}

	// 如果对话没有 Preload Assistant，则手动获取 Assistant
	if chatEntity.AssistantId != nil && chatEntity.Assistant == nil {
		assistantEntity, err = u.assistantService.GetAssistant(c, *chatEntity.AssistantId)
		if err != nil {
			if errors.Is(err, consts.ErrAssistantNotFound) {
				response.Status(http.StatusNotFound).Error(err).Send()
			} else {
				response.Status(http.StatusInternalServerError).Error(err).Send()
			}
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
	cmd := u.redis.Client.Get(c, u.getCacheKey("entity:"+chatEntity.Id.String()))
	_, err = cmd.Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			response.Status(http.StatusInternalServerError).Error(cmd.Err()).Send()

			return
		}
	} else {
		//chatMessageResponse.StreamId = result

		response.Status(http.StatusConflict).Error(consts.ErrChatStreamNotOpen).Data(chatMessageResponse).Send()
		return
	}

	var needStream = true
	var addMessage = true

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

	// 如果是 RoleSystemOverride，更新 Chat，并且也不需要回复
	if request.Role == schema.RoleSystemOverride {
		// 覆盖消息
		chatEntity.Prompt = &request.Message
		err := u.chatService.UpdateChat(c, chatEntity)
		if err != nil {
			response.Error(err).Send()
			return
		}

		addMessage = false
		needStream = false
	}

	// 如果是 RoleHumanLater
	if request.Role == schema.RoleHumanLater {
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
			var addMemory = true

			if request.Role == schema.RoleHuman {
				// 如果对话没有 Assistant，则默认启用记忆
				if chatEntity.Assistant == nil {
					addMemory = true
				} else {
					// 如果禁用了默认 Prompt
					if chatEntity.Assistant.DisableDefaultPrompt {
						// 依旧可以添加记忆
						addMemory = true
					}
				}
			}

			if assistantEntity != nil && assistantEntity.DisableMemory {
				addMemory = false
			}

			if addMemory && lastChatMessage.Role == schema.RoleAssistant {
				u.addMemory(c, userInfo, lastChatMessage.Content, request.Message)
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

	//// 将历史消息分块
	//cmb, err := u.chatService.ChatToBlock(c, chatEntity)
	//if err != nil {
	//	response.Status(http.StatusInternalServerError).Error(err).Send()
	//	return
	//}
	//
	//// 保存分块
	//go func() {
	//	err := u.chatService.SaveBlock(c, cmb)
	//	if err != nil {
	//		u.logger.Sugar.Error(err)
	//	}
	//}()

	var currentAssistant = assistantEntity

	// 如果消息指定了其他助理，则由其他助理来回复
	if request.AssistantId != nil {
		canUse, err := u.assistantService.CanUse(c, userInfo.Token.Sub, *request.AssistantId)
		if err != nil {
			response.Status(http.StatusForbidden).Error(err).Send()
			return
		}

		if !canUse {
			response.Status(http.StatusForbidden).Error(err).Send()
			return
		}

		// 如果这个 assistant 不是用户的
		messageAssistantEntity, err := u.assistantService.GetAssistant(c, *request.AssistantId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		currentAssistant = messageAssistantEntity
	}

	// 检测是否存在知识库，支持临时使用的助理的知识库
	if needStream && currentAssistant != nil && currentAssistant.LibraryId != nil {
		libraryEntity, err := u.libraryService.GetLibrary(c, *currentAssistant.LibraryId)
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
			chunkContent += libraryResult.Content + "\n"
		}

		// 添加知识库消息
		chatMessages = append(chatMessages, entity.ChatMessage{
			ChatId:      chatEntity.Id,
			AssistantId: &currentAssistant.Id,
			Content:     chunkContent,
			Role:        schema.RoleKnowledge,
		})
	}

	// 添加用户发送的消息
	if currentAssistant != nil {
		chatMessages = append(chatMessages, entity.ChatMessage{
			ChatId:      chatEntity.Id,
			AssistantId: &currentAssistant.Id,
			Content:     request.Message,
			Role:        request.Role,
		})
	} else {
		chatMessages = append(chatMessages, entity.ChatMessage{
			ChatId:  chatEntity.Id,
			Content: request.Message,
			Role:    request.Role,
		})
	}

	// TODO: 如果 request.Message 的大小超过了 1mb, 则转换为文件。转换之前应该先判断助理是否存在知识库
	// Update: 其实我也不知道这个要不要做,感觉做了意义也不大
	// 如果存在知识库，则将文件放入知识库中，否则将不处理
	//if len(request.Message) > 1024*1024 {
	//	// 转换为 file
	//
	//}

	if addMessage {
		for _, chatMessage := range chatMessages {
			err = u.cm.CreateChatMessage(c, &chatMessage)
			if err != nil {
				response.Status(http.StatusInternalServerError).Error(err).Send()
				return
			}
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

func (u *ChatController) addMemory(c context.Context, userInfo *schema.User, lastAssistantResponse string, userInput string) {
	// 如果 request.Message 字数大于 200，则跳过
	if len(userInput) > 200 {
		return
	}

	go func() {
		u.logger.Sugar.Info("memory service adding: ", userInput)

		err := u.memoryService.Add(c, userInput, lastAssistantResponse, userInfo.Token.Sub)
		if err != nil {
			u.logger.Sugar.Error("memory service add error: ", err)
		}
	}()
}
