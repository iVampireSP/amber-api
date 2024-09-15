package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

// ListAssistantApiKeys godoc
// @Summary      获取 Assistant API Key列表
// @Description  此 API 可以创建一个 Assistant API Key，可以将你的 Assistant 公开出去使用。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantApiKeyListRequest  path  schema.AssistantApiKeyListRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.AssistantApiKey}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/keys [get]
func (u *AssistantController) ListAssistantApiKeys(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantApiKeyListRequest = schema.AssistantApiKeyListRequest{}
	if err := c.ShouldBindUri(&assistantApiKeyListRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantApiKeyListRequest.AssistantId)
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

	assistantApiKeyList, err := u.assistantService.ListApiKey(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistantApiKeyList).Send()
}

// CreateAssistantApiKey godoc
// @Summary      创建 Assistant API Key
// @Description  此方法将会获取一个 Token，用户将会通过这个 Token 来访问你的 Assistant 并调用工具。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantApiKeyCreateRequest  path  schema.AssistantApiKeyCreateRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.AssistantApiKey}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/keys [post]
func (u *AssistantController) CreateAssistantApiKey(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantApiKeyCreateRequest = schema.AssistantApiKeyCreateRequest{}
	if err := c.ShouldBindUri(&assistantApiKeyCreateRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantApiKeyCreateRequest.AssistantId)
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

	apiKey, err := u.assistantService.CrateApiKey(c, assistantEntity)
	if err != nil {
		return
	}

	response.Status(http.StatusOK).Data(apiKey).Send()
}

// DeleteAssistantApiKey godoc
// @Summary      删除 Assistant API Key
// @Description  此方法将会删除API Key，删除后，API Key将会立即失效。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantApiKeyDeleteRequest  path  schema.AssistantApiKeyDeleteRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/keys/{key_id} [delete]
func (u *AssistantController) DeleteAssistantApiKey(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantApiKeyDeleteRequest = schema.AssistantApiKeyDeleteRequest{}
	if err := c.ShouldBindUri(&assistantApiKeyDeleteRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantApiKeyDeleteRequest.AssistantId)

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

	apiKey, err := u.assistantService.GetApiKey(c, assistantApiKeyDeleteRequest.KeyId)

	if err != nil {
		if errors.Is(err, consts.ErrApiKeyNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
		}
		return
	}

	err = u.assistantService.DeleteApiKey(c, apiKey)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Send()
}

//func (u *AssistantController) UpdateAssistantApiKey(c *gin.Context) {}

func (u *AssistantController) GetAssistantApiKey(c *gin.Context) {
}
