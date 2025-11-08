package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alfian57/golang-todo/pkg/config"
	"github.com/Alfian57/golang-todo/pkg/database"
	"github.com/Alfian57/golang-todo/pkg/logger"
)

// @title Golang Todo API
// @version 1.0
// @description API documentation for Golang Todo
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	log := logger.NewZapLogger(cfg.App.Mode)
	if zapLogger, ok := log.(*logger.ZapLogger); ok {
		defer zapLogger.Sync()
	}

	// Init Swagger Info
	config.InitSwagger()

	// Initialize database connection
	db, err := database.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize database", logger.F("error", err))
	}

	// Setup server
	srv := initServer(db, cfg, log)

	// Start server in goroutine
	go func() {
		log.Info("Starting server",
			logger.F("address", srv.Addr),
			logger.F("mode", cfg.App.Mode),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server",
				logger.F("error", err),
			)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown",
			logger.F("error", err),
		)
	}

	log.Info("Server exited")
}
