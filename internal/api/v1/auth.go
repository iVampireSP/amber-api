package v1

import (
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"

	"github.com/gin-gonic/gin"
)

// AuthController 身份认证控制器
type AuthController struct {
	authService *auth.Service
}

// NewAuthController 创建身份认证控制器
func NewAuthController(authService *auth.Service) *AuthController {
	return &AuthController{authService}
}

// Register godoc
// @Summary      用户注册
// @Description  注册新用户账号
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        registerRequest  body  schema.RegisterRequest  true  "注册信息"
// @Success      200  {object}  schema.ResponseBody{data=schema.TokenResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var response = schema.NewResponse(c)
	var request schema.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	tokenResponse, err := ac.authService.Register(c, &request)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(tokenResponse).Send()
}

// Login godoc
// @Summary      用户登录
// @Description  使用账号密码登录
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body  schema.LoginRequest  true  "登录信息"
// @Success      200  {object}  schema.ResponseBody{data=schema.TokenResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var response = schema.NewResponse(c)
	var request schema.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	tokenResponse, err := ac.authService.Login(c, &request)
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(tokenResponse).Send()
}

// GetCurrentUser godoc
// @Summary      获取当前用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  schema.ResponseBody{data=schema.UserResponse}
// @Failure      401  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/auth/current [get]
func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	var response = schema.NewResponse(c)

	user, err := ac.authService.GetCurrentUser(c)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	response.Status(http.StatusOK).Data(user.ToResponse()).Send()
}
