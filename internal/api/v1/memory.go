package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/memory"
	"rag-new/pkg/consts"
)

type MemoryController struct {
	authService   *auth.Service
	memoryService *memory.Service
	logger        *logger.Logger
	config        *conf.Config
}

func NewMemoryController(authService *auth.Service, memoryService *memory.Service, logger *logger.Logger, conf *conf.Config) *MemoryController {
	return &MemoryController{authService, memoryService, logger, conf}
}

// List godoc
// @Summary      获取所有的记忆
// @Tags         memoires
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  schema.ResponseBody{data=entity.Memory}
// @Failure      400  {object}  schema.ResponseBody{}
// @Router       /api/v1/memories [get]
func (mc *MemoryController) List(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := mc.authService.GetUserId(c)

	memories, err := mc.memoryService.GetMemories(c, userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(memories).Send()
	return
}

// Delete godoc
// @Summary      删除指定的记忆
// @Tags         memoires
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        schema.MemoryDeleteRequest  path  	schema.MemoryDeleteRequest  true  "schema.MemoryDeleteRequest"
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/memories/{id} [delete]
func (mc *MemoryController) Delete(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := mc.authService.GetUserId(c)

	var memoryDeleteRequest schema.MemoryDeleteRequest
	if err := c.ShouldBindUri(&memoryDeleteRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err)
		return
	}

	exists, err := mc.memoryService.Exists(c, uint(memoryDeleteRequest.ID), userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err)
		return
	}

	if !exists {
		response.Status(http.StatusNotFound).Error(consts.ErrMemoryNotFound).Send()
		return
	}

	memories, err := mc.memoryService.GetMemory(c, uint(memoryDeleteRequest.ID), userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	err = mc.memoryService.Delete(c, memories)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
	return
}

// Purge godoc
// @Summary      删除全部记忆
// @Tags         memoires
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  schema.ResponseBody
// @Router       /api/v1/memories/purge [post]
func (mc *MemoryController) Purge(c *gin.Context) {
	var response = schema.NewResponse(c)
	userId := mc.authService.GetUserId(c)

	err := mc.memoryService.Purge(c, userId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusNoContent).Send()
	return
}
