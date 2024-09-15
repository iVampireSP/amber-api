package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

// ListAssistantTokens godoc
// @Summary      获取 Assistant 密钥列表
// @Description  此 API 可以创建一个 Assistant Token，可以将你的 Assistant 对接其他应用。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantTokenListRequest  path  schema.AssistantTokenListRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.AssistantToken}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/tokens [get]
func (u *AssistantController) ListAssistantTokens(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantTokenListRequest = schema.AssistantTokenListRequest{}
	if err := c.ShouldBindUri(&assistantTokenListRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantTokenListRequest.AssistantId)
	if err != nil {
		if errors.Is(err, consts.ErrAssistantNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
		}
		return
	}
	if !u.authService.Compare(c, assistantEntity) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	assistantTokenList, err := u.assistantService.ListToken(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistantTokenList).Send()
}

// CreateAssistantToken godoc
// @Summary      创建 Assistant 密钥
// @Description  此方法将会获取一个 Token，应用将会通过这个 Token 来访问你的 Assistant 并调用工具。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantTokenCreateRequest  path  schema.AssistantTokenCreateRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.AssistantToken}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/tokens [post]
func (u *AssistantController) CreateAssistantToken(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantTokenCreateRequest = schema.AssistantTokenCreateRequest{}
	if err := c.ShouldBindUri(&assistantTokenCreateRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantTokenCreateRequest.AssistantId)
	if err != nil {
		if errors.Is(err, consts.ErrAssistantNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
		}
		return
	}
	if !u.authService.Compare(c, assistantEntity) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	token, err := u.assistantService.CrateToken(c, assistantEntity)
	if err != nil {
		return
	}

	response.Status(http.StatusOK).Data(token).Send()
}

// DeleteAssistantToken godoc
// @Summary      删除 Assistant 密钥
// @Description  此方法将会删除密钥，删除后，密钥将会立即失效。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantTokenDeleteRequest  path  schema.AssistantTokenDeleteRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/tokens/{token_id} [delete]
func (u *AssistantController) DeleteAssistantToken(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantTokenDeleteRequest = schema.AssistantTokenDeleteRequest{}
	if err := c.ShouldBindUri(&assistantTokenDeleteRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantTokenDeleteRequest.AssistantId)

	if err != nil {

		if errors.Is(err, consts.ErrAssistantNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
		}
		return
	}
	if !u.authService.Compare(c, assistantEntity) {

		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	token, err := u.assistantService.GetToken(c, assistantTokenDeleteRequest.TokenId)

	if err != nil {
		if errors.Is(err, consts.ErrTokenNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
		}
		return
	}

	err = u.assistantService.DeleteToken(c, token)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Send()
}

//func (u *AssistantController) UpdateAssistantToken(c *gin.Context) {}

func (u *AssistantController) GetAssistantToken(c *gin.Context) {
}
