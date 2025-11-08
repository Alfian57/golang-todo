package auth

import (
	"github.com/Alfian57/golang-todo/pkg/config"
	"github.com/Alfian57/golang-todo/pkg/logger"
	"github.com/Alfian57/golang-todo/pkg/middleware"
	"github.com/Alfian57/golang-todo/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config, log logger.Logger) {
	jwtUtils := utils.NewJWTUtils(cfg)
	isDebug := cfg.App.Mode != "release"

	authRepository := NewAuthRepository(db)
	authService := NewAuthService(authRepository, jwtUtils, log, isDebug)
	authHandler := NewAuthHandler(authService, log)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/logout", middleware.AuthMiddleware(jwtUtils, isDebug), authHandler.Logout)
	}
}
