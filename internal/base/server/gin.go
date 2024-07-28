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
	swaggerRouter *router.SwaggerRouter,
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
	r.Use(middleware.GinLogger.GinLogger)
	rootGroup := r.Group("")

	rootGroup.GET("/healthz", func(ctx *gin.Context) { ctx.String(200, "OK") })

	// swagger
	swaggerRouter.Register(rootGroup)

	apiV1 := rootGroup.Group("/api/v1")
	{
		apiV1.Use(middleware.JSONResponse.ContentTypeJSON)
		apiV1.Use(middleware.Auth.RequireJWTIDToken)
		apiRouter.InitApiRouter(apiV1)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return r
}
