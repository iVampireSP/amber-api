package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/batch"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/chat_message"
	"rag-new/internal/service/tool"
	"rag-new/pkg/consts"
	"strconv"
)

type AssistantController struct {
	authService        *auth.Service
	toolService        *tool.Service
	assistantService   *assistant.Service
	chatService        *chat.Service
	chatMessageService *chat_message.Service
	batch              *batch.Batch
}

func NewAssistantController(
	authService *auth.Service,
	toolService *tool.Service,
	assistantService *assistant.Service,
	chatService *chat.Service,
	chatMessageService *chat_message.Service,
	batch *batch.Batch,
) *AssistantController {
	return &AssistantController{authService, toolService, assistantService, chatService, chatMessageService, batch}
}

// List godoc
// @Summary      获取 Assistant 列表
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Assistant}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants [get]
func (u *AssistantController) List(c *gin.Context) {
	var response = schema.NewResponse(c)
	assistants, err := u.assistantService.ListAssistantFromUserId(c, u.authService.GetUserId(c))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistants).Send()
}

// GetAssistant godoc
// @Summary      获取指定的 Assistant
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.Assistant}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id} [get]
func (u *AssistantController) GetAssistant(c *gin.Context) {
	var response = schema.NewResponse(c)
	assistantId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, int64(assistantId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if assistantEntity.ID == consts.NoRecord || assistantEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistantEntity).Send()
}

// CreateAssistant godoc
// @Summary      创建 Assistant
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        assistant  body  schema.AssistantCreateRequest  true  "Assistant"
// @Success      200  {object}  schema.ResponseBody{data=entity.Assistant}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants [post]
func (u *AssistantController) CreateAssistant(c *gin.Context) {
	var createReq schema.AssistantCreateRequest
	var response = schema.NewResponse(c)

	if err := c.ShouldBindJSON(&createReq); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	createReq.UserId = u.authService.GetUserId(c)

	assistants, err := u.assistantService.CreateAssistant(c, &createReq)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistants).Send()
}

// DeleteAssistant godoc
// @Summary      删除 Assistant
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Assistant ID"
// @Success      204
// @Failure      500  {object}  schema.ResponseBody{}
// @Failure      404  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id} [delete]
func (u *AssistantController) DeleteAssistant(c *gin.Context) {
	var response = schema.NewResponse(c)

	assistantId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, int64(assistantId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if assistantEntity.ID == consts.NoRecord || assistantEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	// 如果已经绑定过工具，则不能删除
	toolEntity, err := u.assistantService.ListAssistantTool(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if len(toolEntity) > 0 {
		response.Status(http.StatusBadRequest).Error(consts.ErrAssistantHasBindToolCantDelete).Send()
		return
	}

	// batch
	var adb = &batch.AssistantDeleteBatch{
		AssistantEntity:    assistantEntity,
		AssistantService:   u.assistantService,
		ChatService:        u.chatService,
		ChatMessageService: u.chatMessageService,
	}

	err = u.batch.AssistantDelete(c, adb)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
}

// ListTool godoc
// @Summary      获取 Assistant 所绑定的 Tool
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.AssistantToolType}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/tools [get]
func (u *AssistantController) ListTool(c *gin.Context) {
	var response = schema.NewResponse(c)

	assistantId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, int64(assistantId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if assistantEntity.ID == consts.NoRecord || assistantEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	toolList, err := u.assistantService.ListAssistantToolWithType(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(toolList).Send()
}

// BindTool godoc
// @Summary      绑定 Tool
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Assistant ID"
// @Param        tool_id  path  int  true  "Tool ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.AssistantTool}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/tools/{tool_id} [post]
func (u *AssistantController) BindTool(c *gin.Context) {
	var response = schema.NewResponse(c)

	assistantIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	toolIdInt, err := strconv.Atoi(c.Param("tool_id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	assistantId := int64(assistantIdInt)
	toolId := int64(toolIdInt)

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if assistantEntity.ID == consts.NoRecord || assistantEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	toolEntity, err := u.toolService.GetTool(c, toolId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if toolEntity.ID == consts.NoRecord || toolEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrToolNotFound).Send()
		return
	}

	assistantTool, err := u.assistantService.BindTool(c, assistantEntity, toolEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(&assistantTool).Send()
}

// UnbindTool godoc
// @Summary      解绑 Tool
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  int  true  "Assistant ID"
// @Param        tool_id  path  int  true  "Tool ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.AssistantTool}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/tools/{tool_id} [delete]
func (u *AssistantController) UnbindTool(c *gin.Context) {
	var response = schema.NewResponse(c)

	assistantIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	toolIdInt, err := strconv.Atoi(c.Param("tool_id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	assistantId := int64(assistantIdInt)
	toolId := int64(toolIdInt)

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if assistantEntity.ID == consts.NoRecord || assistantEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	toolEntity, err := u.toolService.GetTool(c, toolId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if toolEntity.ID == consts.NoRecord || toolEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrToolNotFound).Send()
		return
	}

	err = u.assistantService.UnbindTool(c, assistantEntity, toolEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Send()
}
