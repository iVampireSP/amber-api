package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"rag-new/docs"
)

type SwaggerRouter struct {
	//config *conf.Config
}

func NewSwaggerRoute() *SwaggerRouter {
	return &SwaggerRouter{}
}

func (a *SwaggerRouter) initSwaggerDocs() {
	docs.SwaggerInfo.Title = "RAG"
	docs.SwaggerInfo.Description = "RAG"
	docs.SwaggerInfo.Version = "v0.0.1"
	docs.SwaggerInfo.BasePath = "/"
}

func (a *SwaggerRouter) Register(r *gin.RouterGroup) {
	a.initSwaggerDocs()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
