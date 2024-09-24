package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/schema"
	_ "rag-new/pkg/page"
)

// AssistantPublicList godoc
// @Summary      获取公开的助理列表
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.ListPublicAssistantReq  query  schema.ListPublicAssistantReq true  "schema.ListPublicAssistantReq"
// @Success      200  {object}  schema.ResponseBody{data=page.PagedResult[schema.AssistantPublic]}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/public [get]
func (u *AssistantController) AssistantPublicList(c *gin.Context) {
	var response = schema.NewResponse(c)
	var listPublicAssistantReq = &schema.ListPublicAssistantReq{}
	if err := c.ShouldBind(listPublicAssistantReq); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	assistants, err := u.assistantService.ListPublicAssistant(c, listPublicAssistantReq)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistants).Send()
}
