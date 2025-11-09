package todo

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

	todoRepository := NewTodoRepository(db)
	todoService := NewTodoService(todoRepository, log, isDebug)
	todoHandler := NewTodoHandler(todoService, log)

	todoGroup := router.Group("/todo", middleware.AuthMiddleware(jwtUtils, isDebug))
	{
		todoGroup.GET("/", todoHandler.GetAll)
		todoGroup.POST("/", todoHandler.Create)
		todoGroup.PUT("/:id", todoHandler.Update)
		todoGroup.DELETE("/:id", todoHandler.Delete)
	}
}
