package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/library"
	"rag-new/pkg/consts"
)

type LibraryController struct {
	libraryService *library.Service
	authService    *auth.Service
}

func NewLibraryController(libraryService *library.Service, authService *auth.Service) *LibraryController {
	return &LibraryController{
		libraryService,
		authService,
	}
}

// List godoc
// @Summary      获取所有的资料库
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Library}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/libraries [get]
func (lc *LibraryController) List(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	libraries, err := lc.libraryService.ListLibraryByUserId(c, userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(libraries).Send()
	return
}

// Delete godoc
// @Summary      删除资料库
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryIdRequest  path  	schema.LibraryIdRequest  true  "schema.LibraryIdRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/libraries/{id} [delete]
func (lc *LibraryController) Delete(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	var libraryIdRequest = &schema.LibraryIdRequest{}
	if err := c.ShouldBindUri(libraryIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	libraryEntity, err := lc.libraryService.GetLibraryByUserId(c, userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	err = lc.libraryService.DeleteLibrary(c, libraryEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
	return
}

// Update godoc
// @Summary      更新资料库
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryIdRequest  path  	schema.LibraryIdRequest  true  "schema.LibraryIdRequest"
// @Param        schema.LibraryUpdateRequest  body  	schema.LibraryUpdateRequest  true  "schema.LibraryUpdateRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/libraries/{id} [patch]
func (lc *LibraryController) Update(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	var libraryIdRequest = &schema.LibraryIdRequest{}
	if err := c.ShouldBindUri(libraryIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}
	var libraryUpdateRequest = &schema.LibraryUpdateRequest{}
	if err := c.ShouldBindJSON(libraryUpdateRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	libraryEntity, err := lc.libraryService.GetLibraryByUserId(c, userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if libraryUpdateRequest.Name != "" {
		libraryEntity.Name = libraryUpdateRequest.Name
	}

	if libraryUpdateRequest.Description != nil {
		libraryEntity.Description = libraryUpdateRequest.Description
	}

	err = lc.libraryService.UpdateLibrary(c, libraryEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
	return
}

// ListDocuments godoc
// @Summary      列出资料库以及资料库下的文档
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryIdRequest  path  	schema.LibraryIdRequest  true  "schema.LibraryIdRequest"
// @Param        schema.LibraryUpdateRequest  body  	schema.LibraryUpdateRequest  true  "schema.LibraryUpdateRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/libraries/{id}/documents [get]
func (lc *LibraryController) ListDocuments(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	var libraryIdRequest = &schema.LibraryIdRequest{}
	if err := c.ShouldBindUri(libraryIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	libraryEntity, err := lc.libraryService.GetLibraryDocumentsById(c, libraryIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if libraryEntity == nil {
		response.Status(http.StatusNoContent).Error(err).Send()
		return
	}

	if libraryEntity.UserId != userId {
		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
		return
	}

	response.Status(http.StatusOK).Data(libraryEntity).Send()
	return
}

func (lc *LibraryController) UpdateDocument(c *gin.Context) {}

// DeleteDocument godoc
// @Summary      删除指定的文档
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryIdRequest  path  	schema.LibraryIdRequest  true  "schema.LibraryIdRequest"
// @Param        schema.LibraryUpdateRequest  body  	schema.LibraryUpdateRequest  true  "schema.LibraryUpdateRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/libraries/{id}/documents/{document_id} [delete]
func (lc *LibraryController) DeleteDocument(c *gin.Context) {}
