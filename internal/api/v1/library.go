package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/library"
)

type LibraryController struct {
	libraryService *library.Service
}

func NewLibraryController(libraryService *library.Service) *LibraryController {
	return &LibraryController{
		libraryService,
	}
}

// List godoc
// @Summary      获取所有的资料库
// @Tags         libraries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  schema.ResponseBody{data=entity.Memory}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/libraries [get]
func (mc *LibraryController) List(c *gin.Context) {
	//var response = schema.NewResponse(c)
	//userId := mc.authService.GetUserId(c)
	//
	//
	//response.Status(http.StatusOK).Data(memories).Send()
	//return
}

// Delete godoc
// @Summary      删除资料库
// @Tags         memoires
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.MemoryDeleteRequest  path  	schema.MemoryDeleteRequest  true  "schema.MemoryDeleteRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/libraries/{id} [delete]
func (mc *LibraryController) Delete(c *gin.Context) {
	var response = schema.NewResponse(c)

	response.Status(http.StatusNoContent).Send()
	return
}
