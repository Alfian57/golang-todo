package main

import (
	"net/http"

	"github.com/Alfian57/golang-todo/internal/auth"
	"github.com/Alfian57/golang-todo/pkg/config"
	"github.com/Alfian57/golang-todo/pkg/logger"
	"github.com/Alfian57/golang-todo/pkg/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func initServer(db *gorm.DB, cfg *config.Config, log logger.Logger) *http.Server {
	r := gin.New()

	// Use custom Zap middleware with injected logger
	r.Use(middleware.ZapRecovery(log))
	r.Use(middleware.ZapLogger(log))

	api := r.Group("api")
	v1 := api.Group("v1")

	// Register Routes
	auth.RegisterRoutes(v1, db, cfg, log)

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    cfg.App.URL,
		Handler: r,
	}

	return srv
}
