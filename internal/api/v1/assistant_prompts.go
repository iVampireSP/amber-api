package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strings"
)

// ListScenePrompt godoc
// @Summary      列出当前助理的场景 Prompt
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        AssistantIDRequest  path  schema.AssistantIDRequest  true  "AssistantIDRequest"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.ScenePrompt}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/scene_prompts [get]
func (u *AssistantController) ListScenePrompt(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantIdRequest = &schema.AssistantIDRequest{}
	if err := c.ShouldBindUri(assistantIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	assistant, err := u.assistantService.GetAssistant(c, assistantIdRequest.ID)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if assistant.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	scenePrompts, err := u.assistantService.GetAssistantScenePrompts(c, assistant)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(scenePrompts).Send()
}

// CreateScenePrompt godoc
// @Summary      创建场景 Prompt
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        AssistantIDRequest  path  schema.AssistantIDRequest  true  "AssistantIDRequest"
// @Param        CreateAssistantScenePromptRequest  body  schema.CreateAssistantScenePromptRequest  true  "CreateAssistantScenePromptRequest"
// @Success      200  {object}  schema.ResponseBody{data=entity.ScenePrompt}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/scene_prompts [post]
func (u *AssistantController) CreateScenePrompt(c *gin.Context) {
	var createReq schema.CreateAssistantScenePromptRequest
	var response = schema.NewResponse(c)

	var assistantIdRequest = &schema.AssistantIDRequest{}
	if err := c.ShouldBindUri(assistantIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	if err := c.ShouldBindJSON(&createReq); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	// 左右去除空格
	createReq.Label = strings.TrimSpace(createReq.Label)

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantIdRequest.ID)
	if err != nil {
		if errors.Is(err, consts.ErrAssistantNotFound) {
			response.Status(http.StatusNotFound).Error(err).Send()
		} else {
			response.Status(http.StatusInternalServerError).Error(err).Send()
		}
		return
	}

	if assistantEntity.UserId != u.authService.GetUserId(c) {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantNotFound).Send()
		return
	}

	// 检测是否存在
	count, err := u.assistantService.CountAssistantScenePrompts(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if count > MaxScenePrompts {
		response.Status(http.StatusBadRequest).Error(consts.ErrAssistantScenePromptMax).Send()
		return
	}

	exists, err := u.assistantService.GetScenePromptExistsByLabel(c, createReq.Label, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if exists {
		response.Status(http.StatusBadRequest).Error(consts.ErrAssistantScenePromptExists).Send()
		return
	}

	scenePrompt, err := u.assistantService.CreateScenePrompt(c, createReq.Label, createReq.Prompt, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(scenePrompt).Send()
}

// DeleteScenePrompt godoc
// @Summary      删除场景 Prompt
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        AssistantIDRequest  path  schema.AssistantIDRequest  true  "AssistantIDRequest"
// @Param        AssistantSceneIdRequest  path  schema.AssistantSceneIdRequest  true  "AssistantSceneIdRequest"
// @Success      204
// @Failure      500  {object}  schema.ResponseBody{}
// @Failure      404  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/{id}/scene_prompts/{scene_id} [delete]
func (u *AssistantController) DeleteScenePrompt(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantIdRequest = &schema.AssistantIDRequest{}
	if err := c.ShouldBindUri(assistantIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var assistantSceneIdRequest = &schema.AssistantSceneIdRequest{}
	if err := c.ShouldBindUri(assistantSceneIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantIdRequest.ID)
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

	scenePrompt, err := u.assistantService.GetScenePromptById(c, assistantSceneIdRequest.SceneId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if scenePrompt.AssistantId != assistantEntity.Id {
		response.Status(http.StatusNotFound).Error(consts.ErrAssistantScenePromptNotFound).Send()
		return
	}

	if err := u.assistantService.DeleteScenePrompt(c, scenePrompt); err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
}
