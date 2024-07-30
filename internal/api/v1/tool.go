package v1

import (
	"net/http"
	_ "rag-new/internal/entity" // 需要保留这行，否则 swag go 解析有问题
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/tool"
	"rag-new/pkg/consts"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ToolController struct {
	toolService *tool.Service
	authService *auth.Service
}

func NewToolController(toolService *tool.Service, authService *auth.Service) *ToolController {
	return &ToolController{
		toolService,
		authService,
	}
}

// List godoc
// @Summary      List Tool
// @Description  List tools
// @Tags         tool
// @Accept       json
// @Produce      json
// @Success      200  {object}  schema.ResponseBody{data=schema.CurrentUserResponse}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/tools [get]
func (t *ToolController) List(c *gin.Context) {
	var response = schema.NewResponse(c)
	tools, err := t.toolService.ListToolFromUserId(c, t.authService.GetUserId(c))
	if err != nil {
		response.Error(err).Send()
		return
	}

	response.Data(tools).Send()
}

// CreateTool godoc
// @Summary      Create Tool
// @Description  Create tool
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        tool  body  schema.ToolCreateRequest  true  "Tool"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Tool}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/tools [post]
func (t *ToolController) CreateTool(c *gin.Context) {
	var response = schema.NewResponse(c)

	// bind req
	var req schema.ToolCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(err).Send()
		return
	}

	exists, err := t.toolService.CheckTool(c, req.Url, t.authService.GetUserId(c))
	if err != nil {
		response.Error(err).Send()
		return
	}

	if exists {
		response.Error(consts.ErrToolAlreadyExists).Send()
		return
	}

	// create
	toolEntity, err := t.toolService.CreateTool(c, &req, t.authService.GetUserId(c))
	if err != nil {
		response.Error(err).Send()
		return
	}

	response.Data(&toolEntity).Send()
}

// GetTool godoc
// @Summary      Get Tool
// @Description  Get tool
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Tool ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.Tool}
// @Failure      400  {object}  schema.ResponseBody{}
// @Failure      404  {object}  schema.ResponseBody{}
// @Router       /api/v1/tools/{id} [get]
func (t *ToolController) GetTool(c *gin.Context) {
	var response = schema.NewResponse(c)
	toolId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	getTool, err := t.toolService.GetTool(c, int64(toolId))
	if err != nil {
		response.Error(err).Send()
		return
	}

	if getTool.ID == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrToolNotFound).Send()
		return
	}

	if getTool.UserId != t.authService.GetUserId(c) {
		response.Error(consts.ErrToolNotYours).Send()
		return
	}

	response.Data(getTool).Send()
}

// DeleteTool godoc
// @Summary      DeleteTool
// @Description  DeleteTool
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Tool ID"
// @Success      200  {object}  nil
// @Failure      400  {object}  schema.ResponseBody{}
// @Failure      404  {object}  schema.ResponseBody{}
// @Router       /api/v1/tools/{id} [delete]
func (t *ToolController) DeleteTool(c *gin.Context) {
	var response = schema.NewResponse(c)
	toolId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(err).Send()
		return
	}

	getTool, err := t.toolService.GetTool(c, int64(toolId))
	if err != nil {
		response.Error(err).Send()
		return
	}

	if getTool.ID == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrToolNotFound).Send()
		return
	}

	if getTool.UserId != t.authService.GetUserId(c) {
		response.Error(consts.ErrToolNotYours).Send()
		return
	}

	err = t.toolService.DeleteTool(c, int64(toolId))
	if err != nil {
		response.Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
}
