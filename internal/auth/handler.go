package auth

import (
	"github.com/Alfian57/golang-todo/common/models"
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

	response := handler.authService.Login(ctx, req)
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

	response := handler.authService.Register(ctx, req)
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
// @Failure      401  {object}  models.Response
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/logout [post]
func (handler AuthHandler) Logout(ctx *gin.Context) {
	var req LogoutRequest
	if !utils.ValidateRequest(ctx, &req) {
		return
	}

	// Get user ID from JWT claims (optional - for logging purposes)
	// This demonstrates how to use GetUserIDFromContext helper
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		handler.log.Warn("Failed to get user ID from context",
			logger.F("operation", "Logout"),
			logger.F("error", err),
		)
		response := utils.UnauthorizedResponse("Unauthorized", err, false)
		ctx.JSON(response.StatusCode, response)
		return
	}

	response := handler.authService.Logout(ctx, req.RefreshToken)
	if response.StatusCode != 201 {
		handler.log.Warn("Logout request failed",
			logger.F("operation", "Logout"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("user_id", userID.String()),
		)
	}

	ctx.JSON(response.StatusCode, response)
}

// @Summary      Get current user
// @Description  Get current authenticated user information
// @Tags         Auth
// @Produce      json
// @Security	 BearerAuth
// @Success      200  {object}  models.Response{data=auth.UserResponse}
// @Failure      401  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/me [get]
func (handler AuthHandler) Me(ctx *gin.Context) {
	// Use helper function to get user ID from JWT claims
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		handler.log.Warn("Failed to get user ID from context",
			logger.F("operation", "Me"),
			logger.F("error", err),
		)

		// Return appropriate error based on error type
		var response models.Response
		switch err {
		case utils.ErrClaimsNotFound:
			response = utils.UnauthorizedResponse("Unauthorized - no token provided", err, false)
		case utils.ErrInvalidClaimsType:
			response = utils.InternalServerErrorResponse("Failed to process token", err, false)
		case utils.ErrInvalidUserID, utils.ErrUserIDNotFound:
			response = utils.UnauthorizedResponse("Invalid token", err, false)
		default:
			response = utils.InternalServerErrorResponse("Failed to process request", err, false)
		}

		ctx.JSON(response.StatusCode, response)
		return
	}

	// Call service to get user details
	response := handler.authService.GetUserByID(ctx, userID)
	if response.StatusCode != 200 {
		handler.log.Warn("Get current user failed",
			logger.F("operation", "Me"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("user_id", userID.String()),
		)
	}

	ctx.JSON(response.StatusCode, response)
}

// @Summary      User refresh token
// @Description  User refresh access token with refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      RefreshTokenRequest  true  "Refresh Token Request"
// @Success      201  {object}  models.Response{data=auth.RefreshTokenResponse}
// @Failure      404  {object}  models.Response
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/refresh-token [post]
func (handler AuthHandler) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if !utils.ValidateRequest(ctx, &req) {
		return
	}

	// Call service with just the refresh token
	// Service will validate the refresh token and return new tokens
	response := handler.authService.RefreshToken(ctx, req.RefreshToken)
	if response.StatusCode != 201 {
		handler.log.Warn("Refresh Token request failed",
			logger.F("operation", "Refresh Token"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
		)
	}

	ctx.JSON(response.StatusCode, response)
}
