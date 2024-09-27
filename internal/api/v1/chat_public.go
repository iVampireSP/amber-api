package v1

import (
	"errors"
	"net/http"
	"rag-new/internal/entity"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CreatePublicChat godoc
// @Summary      通过 API 创建一个公开的对话记录
// @Tags         chat_public
// @Accept       json
// @Produce      json
// @Param        schema.ChatPublicRequest  body  schema.ChatPublicRequest true  "ChatPublicRequest"
// @Success      200  {object}  schema.ResponseBody{data=entity.Chat}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/chat_public [post]
func (u *ChatController) CreatePublicChat(c *gin.Context) {
	var response = schema.NewResponse(c)
	var request schema.ChatPublicRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	assistantShare, err := u.assistantService.GetByKey(c, request.AssistantKey)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var createGuestChatRequest = &schema.ChatGuestCreateRequest{
		Name:        request.Name,
		AssistantId: assistantShare.AssistantId,
		GuestID:     request.GuestId,
	}

	// 创建临时对话
	chat, err := u.chatService.CreateGuestChat(c, createGuestChatRequest)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(chat).Send()
}

// GetChatPublic  godoc
// @Summary      获取公开对话
// @Tags         chat_public
// @Accept       json
// @Produce      json
// @Param        schema.ChatPublicListRequest  body  schema.ChatPublicListRequest true  "ChatPublicListRequest"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Chat}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/chat_public [get]
func (u *ChatController) GetChatPublic(c *gin.Context) {
	var response = schema.NewResponse(c)
	var request schema.ChatPublicListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	chats, err := u.chatService.ListChatFromGuestId(c, request.GuestId)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(chats).Send()
}

// GetPublicChatMessages  godoc
// @Summary      获取公开对话的聊天记录
// @Tags         chat_public
// @Accept       json
// @Produce      json
// @Param        schema.GetPublicChatMessageRequestParams  path  schema.GetPublicChatMessageRequestParams  true  "GetPublicChatMessageRequestParams"
// @Param        schema.GetPublicChatMessageRequest  query  schema.GetPublicChatMessageRequest true  "ChatPublicListRequest"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.ChatMessage}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/chat_public/{chat_id}/messages [get]
func (u *ChatController) GetPublicChatMessages(c *gin.Context) {
	var response = schema.NewResponse(c)

	var getPublicChatMessageRequestParams schema.GetPublicChatMessageRequestParams
	if err := c.ShouldBindUri(&getPublicChatMessageRequestParams); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var getPublicChatMessageRequest schema.GetPublicChatMessageRequest
	if err := c.ShouldBindQuery(&getPublicChatMessageRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}
	// get assistant by token
	assistantShare, err := u.assistantService.GetByKey(c, getPublicChatMessageRequest.AssistantKey)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, getPublicChatMessageRequestParams.ChatId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		}

		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	// 检查 assistant id 是否一致
	if *chatEntity.AssistantId != assistantShare.AssistantId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	if chatEntity.Owner != schema.OwnerGuest || (chatEntity.GuestId != nil && *chatEntity.GuestId != getPublicChatMessageRequest.GuestId) {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	messagesEntity, err := u.cm.GetChatMessage(c, chatEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
			return
		}
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(messagesEntity).Send()
}

// AddPublicChatMessages  godoc
// @Summary      增加公开对话的聊天记录
// @Tags         chat_public
// @Accept       json
// @Produce      json
// @Param        schema.GetPublicChatMessageRequestParams  path  schema.GetPublicChatMessageRequestParams  true  "GetPublicChatMessageRequestParams"
// @Param        schema.AddPublicChatMessageRequest  body  schema.AddPublicChatMessageRequest true  "AddPublicChatMessageRequest"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chat_public/{chat_id}/messages [post]
func (u *ChatController) AddPublicChatMessages(c *gin.Context) {
	var response = schema.NewResponse(c)

	var getPublicChatMessageRequestParams schema.GetPublicChatMessageRequestParams
	if err := c.ShouldBindUri(&getPublicChatMessageRequestParams); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var addPublicChatMessageRequest schema.AddPublicChatMessageRequest
	if err := c.ShouldBindJSON(&addPublicChatMessageRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	if !slices.Contains(allowedChatMessageRoles, addPublicChatMessageRequest.Role) {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatRoleNotAllowed).Send()
		return
	}

	if addPublicChatMessageRequest.Role == schema.RoleFile {
		response.Status(http.StatusBadRequest).Error(consts.ErrRoleCanNotBeFile).Send()
		return
	}

	chatEntity, err := u.chatService.GetChat(c, getPublicChatMessageRequestParams.ChatId)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	if chatEntity.Owner != schema.OwnerGuest || (chatEntity.GuestId != nil && *chatEntity.GuestId != addPublicChatMessageRequest.GuestId) {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, chatEntity.Id)
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}

	var needStream = true
	// 如果不是 human 或者 hide_human，则不需要回复
	if addPublicChatMessageRequest.Role != schema.RoleHuman && addPublicChatMessageRequest.Role != schema.RoleHideHuman {
		// 不需要生成 ID,直接添加
		needStream = false
	}
	// 如果是 RoleSystemOverride，更新 Chat，并且也不需要回复
	if addPublicChatMessageRequest.Role == schema.RoleSystemOverride {
		// 覆盖消息
		chatEntity.Prompt = &addPublicChatMessageRequest.Message
		err := u.chatService.UpdateChat(c, chatEntity)
		if err != nil {
			response.Error(err).Send()
			return
		}

		needStream = false
	}

	// 如果是 RoleHumanLater
	if addPublicChatMessageRequest.Role == schema.RoleHumanLater {
		needStream = false
	}

	// 检测 chat 是否存在缓存，用于判断是否已经打开了对话
	cmd := u.redis.Client.Get(c, u.getCacheKey("entity:"+chatEntity.Id.String()))
	result, err := cmd.Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			response.Status(http.StatusInternalServerError).Error(cmd.Err()).Send()

			return
		}
	} else {
		// 如果存在，则说明用户没有打开对话，直接返回错误
		chatMessageResponse.StreamId = result

		response.Status(http.StatusConflict).Error(consts.ErrChatStreamNotOpen).Data(chatMessageResponse).Send()
		return
	}

	// 生成访客信息
	var publicUser = &schema.UserPublicInfo{
		Name:      "Guest",
		GuestId:   &addPublicChatMessageRequest.GuestId,
		ChatOwner: schema.OwnerUser,
	}

	// 用户打开了会话且没有正在输出的情况，获取最后一条消息
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
			// 检测角色是否是 human
			// 如果两个消息都是 human，则丢弃上一条消息，修改上一条消息的内容
			lastChatMessage.Content = addPublicChatMessageRequest.Message
			err = u.cm.UpdateMessageContent(c, lastChatMessage)
			if err != nil {
				response.Status(http.StatusInternalServerError).Error(err).Send()
				return
			}

			// 如果 stream id 过期了，但 role 还是 entity.RoleHuman ，则说明没有打开 chat stream，重新生成一个 stream id
			randomStreamId, err := u.generateChatStream(c, chatEntity.Id, publicUser, nil)
			if err != nil {
				response.Status(http.StatusInternalServerError).Error(err).Send()
				return
			}
			chatMessageResponse.StreamId = randomStreamId

			response.Status(http.StatusConflict).Error(consts.ErrChatStreamNotOpenAndOverrideMessage).Data(chatMessageResponse).Send()
			return
		}
	}

	var assistantEntity *entity.Assistant
	if chatEntity.AssistantId != nil {
		assistantEntity, err = u.assistantService.GetAssistant(c, *chatEntity.AssistantId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	// 消息写入列表
	var chatMessages []entity.ChatMessage

	// 检测是否存在知识库
	if needStream && assistantEntity != nil && assistantEntity.LibraryId != nil {
		libraryEntity, err := u.libraryService.GetLibrary(c, *assistantEntity.LibraryId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		// 从知识库获取内容，并添加到历史上下文
		libraryResults, err := u.libraryService.SearchLibrary(c, addPublicChatMessageRequest.Message, libraryEntity)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		var chunkContent = ""
		// 将 libraryResults 拼接起来
		for _, libraryResult := range libraryResults {
			chunkContent += libraryResult.Content + "\n"
		}

		if chunkContent == "" {
			chunkContent = consts.LibraryResultEmptyPrompt
		} else {
			chunkContent = consts.LibraryResultPrompt + "\n" + chunkContent
		}

		// 添加知识库消息
		chatMessages = append(chatMessages, entity.ChatMessage{
			ChatId:      chatEntity.Id,
			AssistantId: &assistantEntity.Id,
			Content:     chunkContent,
			Role:        schema.RoleSystem,
		})
	}

	// 添加用户发送的消息
	chatMessages = append(chatMessages, entity.ChatMessage{
		ChatId:  chatEntity.Id,
		Content: addPublicChatMessageRequest.Message,
		Role:    addPublicChatMessageRequest.Role,
	})

	for _, cm := range chatMessages {
		err = u.cm.CreateChatMessage(c, &cm)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	// 如果需要流式输出的情况
	chatMessageResponse.Stream = needStream
	if needStream {
		randomStreamId, err := u.generateChatStream(c, chatEntity.Id, publicUser, nil)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
		chatMessageResponse.StreamId = randomStreamId
	}

	response.Status(http.StatusOK).Data(chatMessageResponse).Send()
}

// ClearPublicChatMessages  godoc
// @Summary      清空公开对话的聊天记录
// @Tags         chat_public
// @Accept       json
// @Produce      json
// @Param        schema.GetPublicChatMessageRequestParams  path  schema.GetPublicChatMessageRequestParams  true  "GetPublicChatMessageRequestParams"
// @Param        schema.GetPublicChatMessageRequest  body  schema.GetPublicChatMessageRequest true  "GetPublicChatMessageRequest"
// @Success      200
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/chat_public/{chat_id}/clear [post]
func (u *ChatController) ClearPublicChatMessages(c *gin.Context) {
	var response = schema.NewResponse(c)

	var getPublicChatMessageRequestParams schema.GetPublicChatMessageRequestParams
	if err := c.ShouldBindUri(&getPublicChatMessageRequestParams); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var getPublicChatMessageRequest schema.GetPublicChatMessageRequest
	if err := c.ShouldBindJSON(&getPublicChatMessageRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}
	// get assistant by token
	assistantShare, err := u.assistantService.GetByKey(c, getPublicChatMessageRequest.AssistantKey)
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
	if *chatEntity.AssistantId != assistantShare.AssistantId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	if chatEntity.Owner != schema.OwnerGuest || (chatEntity.GuestId != nil && *chatEntity.GuestId != getPublicChatMessageRequest.GuestId) {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	err = u.cm.ClearChatMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
}
