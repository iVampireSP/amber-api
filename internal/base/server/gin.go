package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"rag-new/internal/base/conf"
	"rag-new/internal/middleware"
	"rag-new/internal/router"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"time"
)

func promHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

type HttpServer struct {
	Gin           *gin.Engine
	apiRouter     *router.Api
	swaggerRouter *router.SwaggerRouter
	middleware    *middleware.Middleware
}

// NewHTTPServer new http server.
func NewHTTPServer(
	config *conf.Config,
	apiRouter *router.Api,
	swaggerRouter *router.SwaggerRouter,
	middleware *middleware.Middleware,
) *HttpServer {
	if config.Debug.Enabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	return &HttpServer{
		Gin:           r,
		apiRouter:     apiRouter,
		swaggerRouter: swaggerRouter,
		middleware:    middleware,
	}
}

func (hs *HttpServer) AllowAllCors() {
	var config = cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie", "X-Requested-With", "X-Auth-Token", "Authorization"}
	config.MaxAge = 12 * time.Hour

	hs.Gin.Use(cors.New(config))
}

func (hs *HttpServer) BizRouter() *gin.Engine {
	hs.AllowAllCors()

	rootGroup := hs.Gin.Group("")
	rootGroup.Use(hs.middleware.GinLogger.GinLogger)

	// swagger
	hs.swaggerRouter.Register(rootGroup)

	apiV1 := rootGroup.Group("/api/v1")
	{
		//apiV1.Use(corsMiddleWare)
		apiV1.Use(hs.middleware.JSONResponse.ContentTypeJSON)
		apiV1.Use(hs.middleware.Auth.RequireJWTIDToken)
		hs.apiRouter.InitApiRouter(apiV1)
	}

	apiV1NoAuth := rootGroup.Group("/api/v1")
	{
		//apiV1.Use(corsMiddleWare)
		hs.apiRouter.InitNoAuthApiRouter(apiV1NoAuth)
	}

	hs.Gin.NoRoute(func(ctx *gin.Context) {
		schema.NewResponse(ctx).Status(http.StatusNotFound).Error(consts.ErrPageNotFound).Send()
	})

	return hs.Gin
}

func (hs *HttpServer) MetricRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	metricGroup := r.Group("")
	// prometheus
	metricGroup.GET("/metrics", promHandler(promhttp.Handler()))
	metricGroup.GET("/healthz", func(ctx *gin.Context) { ctx.String(200, "OK") })

	return r
}
