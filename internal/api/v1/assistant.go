package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/tool"
	"rag-new/pkg/consts"
	"strconv"
)

type AssistantController struct {
	authService      *auth.Service
	toolService      *tool.Service
	assistantService *assistant.Service
}

func NewAssistantController(
	authService *auth.Service,
	toolService *tool.Service,
	assistantService *assistant.Service,
) *AssistantController {
	return &AssistantController{authService, toolService, assistantService}
}

// List godoc
// @Summary      获取 Assistant 列表
// @Description  get string by ID
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Assistant}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants [get]
func (u *AssistantController) List(c *gin.Context) {
	assistants, err := u.assistantService.ListAssistantFromUserId(c, u.authService.GetUserId(c))
	if err != nil {
		schema.NewResponse(c).Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	schema.NewResponse(c).Status(http.StatusOK).Data(assistants).Send()
}

// CreateAssistant godoc
// @Summary      创建 Assistant
// @Description  get string by ID
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Param        assistant  body  schema.AssistantCreateRequest  true  "Assistant"
// @Success      200  {object}  schema.ResponseBody{data=entity.Assistant}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants [post]
func (u *AssistantController) CreateAssistant(c *gin.Context) {
	var createReq schema.AssistantCreateRequest

	if err := c.ShouldBindJSON(&createReq); err != nil {
		schema.NewResponse(c).Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	createReq.UserId = u.authService.GetUserId(c)

	assistants, err := u.assistantService.CreateAssistant(c, &createReq)
	if err != nil {
		schema.NewResponse(c).Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	schema.NewResponse(c).Status(http.StatusOK).Data(assistants).Send()
}

// ListTool godoc
// @Summary      获取 Assistant 所绑定的 Tool
// @Description  get string by ID
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.AssistantTool}
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

	toolList, err := u.assistantService.ListAssistantTool(c, int64(assistantId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(toolList).Send()
}

// BindTool godoc
// @Summary      绑定 Tool
// @Description  get string by ID
// @Tags         assistant
// @Accept       json
// @Produce      json
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

	assistantTool, err := u.assistantService.BindTool(c, assistantId, toolId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(&assistantTool).Send()
}
