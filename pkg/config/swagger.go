package config

import (
	"os"

	"github.com/Alfian57/golang-todo/docs"
)

func InitSwagger() {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "golang-todo"
	}

	appUrl := os.Getenv("APP_URL")
	if appUrl == "" {
		appUrl = "localhost:8080"
	}

	docs.SwaggerInfo.Title = appName + " API"
	docs.SwaggerInfo.Description = "API documentation for " + appName
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = appUrl
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
