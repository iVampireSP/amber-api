package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/entity"
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

// CreateLibrary godoc
// @Summary      创建一个资料库
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryCreateRequest  body  	schema.LibraryCreateRequest  true  "schema.LibraryCreateRequest"
// @Success      201  {object}  schema.ResponseBody{data=entity.Library}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/libraries [post]
func (lc *LibraryController) CreateLibrary(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	var libraryCreateRequest = &schema.LibraryCreateRequest{}
	if err := c.ShouldBindJSON(libraryCreateRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var libraryEntity = &entity.Library{
		Name:        libraryCreateRequest.Name,
		Description: libraryCreateRequest.Description,
		UserId:      userId,
	}

	err := lc.libraryService.CreateLibrary(c, libraryEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusCreated).Data(libraryEntity).Send()
	return
}

// GetLibrary godoc
// @Summary      获取一个资料库
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryIdRequest  path  	schema.LibraryIdRequest  true  "schema.LibraryIdRequest"
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Library}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/libraries/{id} [get]
func (lc *LibraryController) GetLibrary(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	var libraryIdRequest = &schema.LibraryIdRequest{}
	if err := c.ShouldBindUri(libraryIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	libraryEntity, err := lc.libraryService.GetLibrary(c, libraryIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if libraryEntity.UserId != userId {
		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
		return
	}

	response.Status(http.StatusOK).Data(libraryEntity).Send()
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

	libraryEntity, err := lc.libraryService.GetLibrary(c, libraryIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if libraryEntity.UserId != userId {
		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
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

	libraryEntity, err := lc.libraryService.GetLibrary(c, libraryIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if libraryEntity.UserId != userId {
		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
		return
	}

	if libraryUpdateRequest.Name != "" {
		libraryEntity.Name = libraryUpdateRequest.Name
	}

	if libraryUpdateRequest.Description != nil {
		libraryEntity.Description = libraryUpdateRequest.Description
	}
	if libraryUpdateRequest.Default != nil {
		libraryEntity.Default = *libraryUpdateRequest.Default
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
// @Success      200  {object}  schema.ResponseBody{data=[]entity.Document}
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

	libraryEntity, err := lc.libraryService.GetLibrary(c, libraryIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if libraryEntity.UserId != userId {
		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
		return
	}

	documents, err := lc.libraryService.GetLibraryDocumentsById(c, libraryIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(documents).Send()
	return
}

// CreateDocument godoc
// @Summary      创建文档
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryIdRequest  path  	schema.LibraryIdRequest  true  "schema.LibraryIdRequest"
// @Param        schema.DocumentCreateRequest  body  	schema.DocumentCreateRequest  true  "schema.DocumentCreateRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/libraries/{id}/documents [post]
func (lc *LibraryController) CreateDocument(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	var libraryIdRequest = &schema.LibraryIdRequest{}
	if err := c.ShouldBindUri(libraryIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	var documentCreateRequest = &schema.DocumentCreateRequest{}
	if err := c.ShouldBindJSON(documentCreateRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	libraryEntity, err := lc.libraryService.GetLibrary(c, libraryIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
	if libraryEntity.UserId != userId {
		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
		return
	}

	documentEntity, err := lc.libraryService.CreateDocument(c, libraryEntity,
		documentCreateRequest.Name, documentCreateRequest.Content)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusCreated).Data(documentEntity).Send()
	return
}

//// UpdateDocument godoc
//// @Summary      更新文档
//// @Tags         libraries
//// @Accept       json
//// @Produce      json
//// @Security     ApiKeyAuth
//// @Param        schema.LibraryAndDocumentIdRequest  path  	schema.LibraryAndDocumentIdRequest  true  "schema.LibraryAndDocumentIdRequest"
//// @Param        schema.LibraryUpdateRequest  body  	schema.LibraryUpdateRequest  true  "schema.LibraryUpdateRequest"
//// @Success      204
//// @Failure      400  {object}  schema.ResponseBody
//// @Router       /api/v1/libraries/{id}/documents/{document_id} [patch]
//func (lc *LibraryController) UpdateDocument(c *gin.Context) {
//	var response = schema.NewResponse(c)
//	userId := lc.authService.GetUserId(c)
//
//	var libraryAndDocumentIdRequest = &schema.LibraryAndDocumentIdRequest{}
//	if err := c.ShouldBindUri(libraryAndDocumentIdRequest); err != nil {
//		response.Status(http.StatusBadRequest).Error(err).Send()
//		return
//	}
//
//	var libraryUpdateRequest = &schema.LibraryUpdateRequest{}
//	if err := c.ShouldBindJSON(libraryUpdateRequest); err != nil {
//		response.Status(http.StatusBadRequest).Error(err).Send()
//		return
//	}
//
//	libraryEntity, err := lc.libraryService.GetLibrary(c, libraryAndDocumentIdRequest.Id)
//	if err != nil {
//		response.Status(http.StatusInternalServerError).Error(err).Send()
//		return
//	}
//
//	if libraryEntity.UserId != userId {
//		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
//		return
//	}
//
//	documentEntity, err := lc.libraryService.GetDocumentFromLibrary(c, libraryEntity,
//		libraryAndDocumentIdRequest.DocumentId)
//	if err != nil {
//		response.Status(http.StatusInternalServerError).Error(err).Send()
//		return
//	}
//
//	documentEntity.Name = libraryUpdateRequest.Name
//	//documentEntity.Description = libraryUpdateRequest.Description
//
//}

// DeleteDocument godoc
// @Summary      删除指定的文档
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.LibraryAndDocumentIdRequest  path  	schema.LibraryAndDocumentIdRequest  true  "schema.LibraryAndDocumentIdRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/libraries/{id}/documents/{document_id} [delete]
func (lc *LibraryController) DeleteDocument(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := lc.authService.GetUserId(c)

	var libraryAndDocumentIdRequest = &schema.LibraryAndDocumentIdRequest{}
	if err := c.ShouldBindUri(libraryAndDocumentIdRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	libraryEntity, err := lc.libraryService.GetLibrary(c, libraryAndDocumentIdRequest.Id)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if libraryEntity.UserId != userId {
		response.Status(http.StatusForbidden).Error(consts.ErrPermissionDenied).Send()
		return
	}

	documentEntity, err := lc.libraryService.GetDocumentFromLibrary(c, libraryEntity,
		libraryAndDocumentIdRequest.DocumentId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// 删除
	err = lc.libraryService.DeleteDocument(c, documentEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
}
