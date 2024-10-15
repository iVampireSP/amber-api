package v1

import (
	"net/http"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"

	"github.com/gin-gonic/gin"
	_ "github.com/iVampireSP/pkg/page"
)

// AssistantPublicList godoc
// @Summary      获取公开的助理列表
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.PaginationRequest  query  schema.PaginationRequest true  "schema.PaginationRequest"
// @Success      200  {object}  schema.ResponseBody{data=page.PagedResult[schema.AssistantPublic]}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/public [get]
func (u *AssistantController) AssistantPublicList(c *gin.Context) {
	var response = schema.NewResponse(c)
	var paginationRequest = &schema.PaginationRequest{}
	if err := c.ShouldBind(paginationRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	assistants, err := u.assistantService.ListPublicAssistant(c, paginationRequest.Page)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistants).Send()
}

// FavoriteAssistant godoc
// @Summary      收藏助理
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        AssistantIDRequest  path  schema.AssistantIDRequest  true  "AssistantIDRequest"
// @Success      200  {object}  schema.ResponseBody{data=schema.AssistantPublic}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/public/{id} [post]
func (u *AssistantController) FavoriteAssistant(c *gin.Context) {
	var response = schema.NewResponse(c)
	var assistantIdRequest = &schema.AssistantIDRequest{}
	if err := c.ShouldBindUri(assistantIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var userId = u.authService.GetUserId(c)

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantIdRequest.ID)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if assistantEntity.UserId == userId {
		response.Status(http.StatusBadRequest).Error(consts.ErrAssistantCannotFavoriteSelf).Send()
		return
	}

	err = u.assistantService.FavoriteAssistant(c, userId, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistantEntity.ToPublic()).Send()
}

// UnFavoriteAssistant godoc
// @Summary      取消收藏助理
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        AssistantIDRequest  path  schema.AssistantIDRequest  true  "AssistantIDRequest"
// @Success      200  {object}  schema.ResponseBody{data=schema.AssistantPublic}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/public/{id} [delete]
func (u *AssistantController) UnFavoriteAssistant(c *gin.Context) {
	var response = schema.NewResponse(c)

	var assistantIdRequest = &schema.AssistantIDRequest{}
	if err := c.ShouldBindUri(assistantIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var userId = u.authService.GetUserId(c)

	assistantEntity, err := u.assistantService.GetAssistant(c, assistantIdRequest.ID)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	err = u.assistantService.UnFavoriteAssistant(c, userId, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistantEntity.ToPublic()).Send()
}

// FavoriteAssistants godoc
// @Summary      收藏的助理列表
// @Tags         assistant
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.PaginationRequest  query  schema.PaginationRequest true  "schema.PaginationRequest"
// @Success      200  {object}  schema.ResponseBody{data=page.PagedResult[schema.AssistantPublic]}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/assistants/favorites [get]
func (u *AssistantController) FavoriteAssistants(c *gin.Context) {
	var response = schema.NewResponse(c)

	var paginationRequest = &schema.PaginationRequest{}
	if err := c.ShouldBind(paginationRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var userId = u.authService.GetUserId(c)

	assistantPublic, err := u.assistantService.ListUserFavoriteAssistants(c, userId, paginationRequest.Page)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(assistantPublic).Send()
}
