package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/tool"
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
// @Failure      400  {object}  schema.ResponseBody{data=schema.EmptyData}
// @Router       /api/v1/assistants [get]
func (u *AssistantController) List(c *gin.Context) {
	assistants, err := u.assistantService.ListAssistantFromUserId(c, u.authService.GetUserId(c))
	if err != nil {
		schema.NewResponse(c).Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	schema.NewResponse(c).Status(http.StatusOK).Data(assistants).Send()
}
