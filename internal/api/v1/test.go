package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
)

type UserController struct {
	authService *auth.Service
}

func NewUserController(authService *auth.Service) *UserController {
	return &UserController{authService}
}

// Test godoc
// @Summary      Greet
// @Description  get string by ID
// @Tags         ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  schema.ResponseBody{data=schema.CurrentUserResponse}
// @Failure      400  {object}  schema.ResponseBody{data=schema.EmptyData}
// @Router       /api/v1/ping [get]
func (u *UserController) Test(c *gin.Context) {
	user := u.authService.GetUser(c)

	var currentUserResponse = &schema.CurrentUserResponse{
		IP:        c.ClientIP(),
		Valid:     user.Valid,
		UserEmail: user.Token.Email,
		UserId:    user.Token.Sub,
	}

	schema.NewResponse(c).Status(http.StatusOK).Data(currentUserResponse).Send()
}
