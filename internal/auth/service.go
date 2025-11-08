package auth

import (
	"context"
	"time"

	"github.com/Alfian57/golang-todo/common/models"
	"github.com/Alfian57/golang-todo/pkg/logger"
	"github.com/Alfian57/golang-todo/pkg/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	authRepository AuthRepository
	jwtUtils       *utils.JWTUtils
	log            logger.Logger
	isDebug        bool
}

func NewAuthService(authRepository AuthRepository, jwtUtils *utils.JWTUtils, log logger.Logger, isDebug bool) AuthService {
	return AuthService{
		authRepository: authRepository,
		jwtUtils:       jwtUtils,
		log:            log,
		isDebug:        isDebug,
	}
}

func (service AuthService) Login(ctx context.Context, username string, password string) models.Response {
	user, err := service.authRepository.FindUserByUsername(ctx, username)
	if err != nil {
		service.log.Debug("User not found during login",
			logger.F("username", username),
			logger.F("error", err),
		)
		return utils.UnauthorizedResponse("Username or password wrong", err, service.isDebug)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		service.log.Debug("Invalid password attempt",
			logger.F("username", username),
		)
		return utils.UnauthorizedResponse("Username or password wrong", nil, service.isDebug)
	}

	accessToken, err := service.jwtUtils.CreateJWT(user.ID.String())
	if err != nil {
		service.log.Error("Failed to create access token",
			logger.F("operation", "Create access token"),
			logger.F("user_id", user.ID.String()),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to create access token", err, service.isDebug)
	}

	refreshToken, err := utils.CreateRefreshToken()
	if err != nil {
		service.log.Error("Failed to create refresh token",
			logger.F("operation", "Create refresh token"),
			logger.F("user_id", user.ID.String()),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to create refresh token", err, service.isDebug)
	}

	_, err = service.authRepository.CreateRefreshToken(ctx, refreshToken, user.ID, service.jwtUtils.GetJWTTTL())
	if err != nil {
		service.log.Error("Failed to save refresh token",
			logger.F("operation", "Save refresh token"),
			logger.F("user_id", user.ID.String()),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to save refresh token", err, service.isDebug)
	}

	responseData := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:        user.ID.String(),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}
	return utils.OkResponse("Success to login", responseData)
}

func (service AuthService) Register(ctx context.Context, username string, password string) models.Response {
	existingUser, err := service.authRepository.FindUserByUsername(ctx, username)
	if err == nil && existingUser.ID != uuid.Nil {
		service.log.Debug("Username already exists during registration",
			logger.F("username", username),
		)
		return utils.UnprocessableEntityResponse("Username already exists", nil, service.isDebug)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		service.log.Error("Failed to hash password",
			logger.F("operation", "Hash password"),
			logger.F("username", username),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to hash password", err, service.isDebug)
	}

	newuser, err := service.authRepository.CreateUser(ctx, username, hashedPassword)
	if err != nil {
		service.log.Error("Failed to create user",
			logger.F("operation", "Create user"),
			logger.F("username", username),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to create user", err, service.isDebug)
	}

	responseData := RegisterResponse{
		User: UserResponse{
			ID:        newuser.ID.String(),
			Username:  newuser.Username,
			CreatedAt: newuser.CreatedAt.Format(time.RFC3339),
		},
	}
	return utils.CreatedResponse("Success to register", responseData)
}

func (service AuthService) Logout(ctx context.Context, token string) models.Response {
	rowsAffected, err := service.authRepository.DeleteRefreshTokenByToken(ctx, token)
	if err != nil {
		service.log.Error("Failed to delete refresh token during logout",
			logger.F("operation", "Logout - delete refresh token"),
			logger.F("refresh_token", token),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to logout", err, service.isDebug)
	}
	if rowsAffected == 0 {
		service.log.Debug("Refresh token not found during logout",
			logger.F("refresh_token", token),
		)
		return utils.NotFoundResponse("Refresh token not exist", nil, service.isDebug)
	}

	return utils.CreatedResponse("Success to logout", nil)
}
