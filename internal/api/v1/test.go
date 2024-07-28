package v1

import (
	"github.com/gin-gonic/gin"
	"rag-new/internal/service/auth"
)

type UserController struct {
	authService *auth.Service
}

func NewUserController(authService *auth.Service) *UserController {
	return &UserController{authService}
}

func (u *UserController) Test(c *gin.Context) {
	user := u.authService.GetUser(c)
	c.JSON(200, gin.H{
		"message": "pong, " + user.Token.Name + ", " + user.Token.Sub,
	})
}
