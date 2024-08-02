package middleware

import "github.com/gin-gonic/gin"

type CORSMiddleware struct{}

func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (*CORSMiddleware) AllowAllCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Max-Age", "1800")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
