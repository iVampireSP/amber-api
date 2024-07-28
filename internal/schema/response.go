package schema

import "github.com/gin-gonic/gin"

type ResponseBody struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

func ResponseMessage(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, &ResponseBody{
		Message: message,
		Code:    code,
		Data:    data,
	})
	c.Abort()
}

func ResponseError(c *gin.Context, code int, err error) {
	c.JSON(code, &ResponseBody{
		Error: err.Error(),
		Code:  code,
	})
	c.Abort()
}

func Response(c *gin.Context, code int, data interface{}) {
	c.JSON(code, &ResponseBody{
		Code: code,
		Data: data,
	})
	c.Abort()
}
