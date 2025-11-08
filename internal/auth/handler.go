package auth

import (
	"github.com/Alfian57/golang-todo/pkg/logger"
	"github.com/Alfian57/golang-todo/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
	log         logger.Logger
}

func NewAuthHandler(authService AuthService, log logger.Logger) AuthHandler {
	return AuthHandler{
		authService: authService,
		log:         log,
	}
}

// @Summary      User login
// @Description  User login with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Login Request"
// @Success      200  {object}  models.Response{data=auth.LoginResponse}
// @Failure      401  {object}  models.Response
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/login [post]
func (handler AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if !utils.ValidateRequest(ctx, &req) {
		return
	}

	response := handler.authService.Login(ctx, req.Username, req.Password)
	if response.StatusCode != 200 {
		handler.log.Warn("Login request failed",
			logger.F("operation", "Login"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("username", req.Username),
		)
	}

	ctx.JSON(response.StatusCode, response)
}

// @Summary      User registration
// @Description  User registration with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Register Request"
// @Success      200  {object}  models.Response{data=auth.RegisterResponse}
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/register [post]
func (handler AuthHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if !utils.ValidateRequest(ctx, &req) {
		return
	}

	response := handler.authService.Register(ctx, req.Username, req.Password)
	if response.StatusCode != 201 {
		handler.log.Warn("Register request failed",
			logger.F("operation", "Register"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("username", req.Username),
		)
	}

	ctx.JSON(response.StatusCode, response)
}

// @Summary      User logout
// @Description  User logout with refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      LogoutRequest  true  "Logout Request"
// @Security	 BearerAuth
// @Success      200  {object}  models.Response
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/logout [post]
func (handler AuthHandler) Logout(ctx *gin.Context) {
	var req LogoutRequest
	if !utils.ValidateRequest(ctx, &req) {
		return
	}

	response := handler.authService.Logout(ctx, req.RefreshToken)
	if response.StatusCode != 200 {
		handler.log.Warn("Logout request failed",
			logger.F("operation", "Logout"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
		)
	}

	ctx.JSON(response.StatusCode, response)
}
