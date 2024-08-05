package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

// ListAssistantShares godoc
// @Summary      获取 Assistant 共享列表
// @Description  此 API 可以创建一个 Assistant 共享 Token，可以将你的 Assistant 公开出去使用。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantShareListRequest  path  schema.AssistantShareListRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.AssistantShare}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/shares [get]
func (u *AssistantController) ListAssistantShares(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantShareListRequest = schema.AssistantShareListRequest{}
	if err := c.ShouldBindUri(&assistantShareListRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantShareListRequest.AssistantId)
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

	assistantShareList, err := u.assistantService.ListShare(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistantShareList).Send()
}

// CreateAssistantShare godoc
// @Summary      创建 Assistant 共享
// @Description  此方法将会获取一个 Token，用户将会通过这个 Token 来访问你的 Assistant 并调用工具。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantShareCreateRequest  path  schema.AssistantShareCreateRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{data=entity.AssistantShare}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/shares [post]
func (u *AssistantController) CreateAssistantShare(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantShareCreateRequest = schema.AssistantShareCreateRequest{}
	if err := c.ShouldBindUri(&assistantShareCreateRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantShareCreateRequest.AssistantId)
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

	share, err := u.assistantService.CrateShare(c, assistantEntity)
	if err != nil {
		return
	}

	response.Status(http.StatusOK).Data(share).Send()
}

// DeleteAssistantShare godoc
// @Summary      删除 Assistant 共享
// @Description  此方法将会删除共享，删除后，共享将会立即失效。
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.AssistantShareDeleteRequest  path  schema.AssistantShareDeleteRequest true  "Assistant ID"
// @Success      200  {object}  schema.ResponseBody{}
// @Failure      500  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/shares/{share_id} [delete]
func (u *AssistantController) DeleteAssistantShare(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantShareDeleteRequest = schema.AssistantShareDeleteRequest{}
	if err := c.ShouldBindUri(&assistantShareDeleteRequest); err != nil {
		response.Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantShareDeleteRequest.AssistantId)

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

	share, err := u.assistantService.GetShare(c, assistantShareDeleteRequest.ShareId)
	fmt.Println(err)

	if err != nil {
		if errors.Is(err, consts.ErrShareNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
		}
		return
	}

	err = u.assistantService.DeleteShare(c, share.ToAssistantShare())
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Send()
}

//func (u *AssistantController) UpdateAssistantShare(c *gin.Context) {}

func (u *AssistantController) GetAssistantShare(c *gin.Context) {
}
