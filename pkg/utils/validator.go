package utils

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateRequest(ctx *gin.Context, req any) bool {
	isDebug := os.Getenv("GIN_MODE") == "debug"

	if err := ctx.ShouldBindJSON(req); err != nil {
		response := UnprocessableEntityResponse("Invalid request body", err, isDebug)
		ctx.JSON(response.StatusCode, response)
		return false
	}

	if err := validate.Struct(req); err != nil {
		response := UnprocessableEntityResponse("Validation failed", err, isDebug)
		ctx.JSON(response.StatusCode, response)
		return false
	}

	return true
}
