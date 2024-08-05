package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"rag-new/internal/entity"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strconv"
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

	assistantShare, err := u.assistantService.GetShareByToken(c, request.AssistantToken)
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
// @Param        schema.GetPublicChatMessageRequest  body  schema.GetPublicChatMessageRequest true  "ChatPublicListRequest"
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
	if err := c.ShouldBindJSON(&getPublicChatMessageRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}
	// get assistant by token
	assistantShare, err := u.assistantService.GetShareByToken(c, getPublicChatMessageRequest.AssistantToken)
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
	if chatEntity.AssistantId != assistantShare.AssistantId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	if chatEntity.Owner != schema.OwnerGuest || chatEntity.GuestId != getPublicChatMessageRequest.GuestId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	messagesEntity, err := u.cm.GetChatMessage(c, chatEntity)
	if err != nil {
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

	chatEntity, err := u.chatService.GetChat(c, getPublicChatMessageRequestParams.ChatId)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	if chatEntity.Owner != schema.OwnerGuest || chatEntity.GuestId != addPublicChatMessageRequest.GuestId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	// 检查状态是否是回复中
	isStreaming, err := u.isStreaming(c, chatEntity.ID)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if isStreaming {
		response.Status(http.StatusBadRequest).Error(consts.ErrChatStreaming).Send()
		return
	}

	var chatMessageResponse = &schema.ChatMessageResponse{}

	var chatIdStr = strconv.Itoa(int(chatEntity.ID))

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

	// last chat message
	lastChatMessage, err := u.cm.GetLatestMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var publicUser = &schema.UserPublicInfo{
		Name:      "Guest",
		Id:        addPublicChatMessageRequest.GuestId,
		ChatOwner: schema.OwnerUser,
	}

	if lastChatMessage.Role == schema.RoleHuman {
		lastChatMessage.Content = addPublicChatMessageRequest.Message
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
	chatMessage.Content = addPublicChatMessageRequest.Message
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
	assistantShare, err := u.assistantService.GetShareByToken(c, getPublicChatMessageRequest.AssistantToken)
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
	if chatEntity.AssistantId != assistantShare.AssistantId {
		response.Status(http.StatusForbidden).Error(err).Send()
		return
	}

	if chatEntity.Owner != schema.OwnerGuest || chatEntity.GuestId != getPublicChatMessageRequest.GuestId {
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
