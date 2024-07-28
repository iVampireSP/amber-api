package schema

import "github.com/gin-gonic/gin"

type ResponseBody struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    any    `json:"data"`
}

type HttpResponse struct {
	body       *ResponseBody
	httpStatus int
	ctx        *gin.Context
}

func NewResponse(c *gin.Context) *HttpResponse {
	return &HttpResponse{
		body:       &ResponseBody{},
		httpStatus: 0,
		ctx:        c,
	}
}

func (r *HttpResponse) Message(message string) *HttpResponse {
	r.body.Message = message

	return r
}

func (r *HttpResponse) Data(data any) *HttpResponse {
	r.body.Data = data
	return r

}

func (r *HttpResponse) Error(err error) *HttpResponse {
	r.body.Error = err.Error()
	return r

}

func (r *HttpResponse) Status(status int) *HttpResponse {
	r.httpStatus = status
	return r

}

func (r *HttpResponse) Send() *HttpResponse {
	r.ctx.JSON(r.httpStatus, r.body)
	return r

}

func (r *HttpResponse) Abort() {
	r.ctx.Abort()
}

//
//func ResponseMessage(c *gin.Context, code int, message string, data interface{}) {
//	c.JSON(code, &ResponseBody{
//		Message: message,
//		Data:    data,
//	})
//	c.Abort()
//}
//
//func ResponseError(c *gin.Context, code int, err error) {
//	c.JSON(code, &ResponseBody{
//		Error: err.Error(),
//	})
//	c.Abort()
//}
//
//func Response(c *gin.Context, code int, data interface{}) {
//	c.JSON(code, &ResponseBody{
//		Data: data,
//	})
//	c.Abort()
//}
