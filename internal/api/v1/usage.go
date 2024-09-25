package v1

import (
	"github.com/gin-gonic/gin"
	"rag-new/internal/schema"
	"rag-new/internal/service/token_usage"
)

type UsageController struct {
	tokenUsage *token_usage.Service
}

func NewUsageController(tokenUsage *token_usage.Service) *UsageController {
	return &UsageController{
		tokenUsage,
	}
}

// GetUsage godoc
// @Summary      获取站点 Usage
// @Tags         usage
// @Accept       json
// @Produce      json
// @Success      200  {object} schema.ResponseBody{data=schema.SiteUsageResponse}
// @Router       /api/v1/usage [get]
func (uc *UsageController) GetUsage(ctx *gin.Context) {
	var response = schema.NewResponse(ctx)

	var siteUsageResponse = schema.SiteUsageResponse{
		MonthTokens:    uc.tokenUsage.GetMonthTokenUsage(ctx),
		MonthToolCalls: uc.tokenUsage.GetMonthToolCallTimes(ctx),
	}

	response.Data(siteUsageResponse).Send()
}
