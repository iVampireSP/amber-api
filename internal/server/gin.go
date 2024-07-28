package server

import (
	"github.com/gin-gonic/gin"
	"rag-new/internal/base/conf"
	"rag-new/internal/middleware"
	"rag-new/internal/router"
)

// NewHTTPServer new http server.
func NewHTTPServer(
	apiRouter *router.Api,
	config *conf.Config,
	middleware *middleware.Middleware,
) *gin.Engine {
	if config.Debug.Enabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.GinLogger.GinLogger())

	r.GET("/healthz", func(ctx *gin.Context) { ctx.String(200, "OK") })

	// The route must be available without logging in
	apiV1 := r.Group("/api/v1")
	apiRouter.InitApiRouter(apiV1)

	return r
}
