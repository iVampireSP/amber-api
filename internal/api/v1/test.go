package v1

import "github.com/gin-gonic/gin"

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
