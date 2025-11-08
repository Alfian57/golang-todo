package utils

import (
	"log"
	"net/http"

	"github.com/Alfian57/golang-todo/common/models"
)

// SuccessResponse creates a success response
func SuccessResponse(statusCode int, message string, data any) models.Response {
	return models.Response{
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
	}
}

// ErrorResponse creates an error response
// If isDebug is true, the error message will be included in the response
func ErrorResponse(statusCode int, message string, err error, isDebug bool) models.Response {
	if err != nil {
		log.Printf("Error: %v", err)
		if isDebug {
			message = err.Error()
		}
	}
	return models.Response{
		Message:    message,
		StatusCode: statusCode,
	}
}

// OkResponse creates a 200 OK response
func OkResponse(message string, data any) models.Response {
	return SuccessResponse(http.StatusOK, message, data)
}

// CreatedResponse creates a 201 Created response
func CreatedResponse(message string, data any) models.Response {
	return SuccessResponse(http.StatusCreated, message, data)
}

// UnauthorizedResponse creates a 401 Unauthorized response
func UnauthorizedResponse(message string, err error, isDebug bool) models.Response {
	return ErrorResponse(http.StatusUnauthorized, message, err, isDebug)
}

// NotFoundResponse creates a 404 Not Found response
func NotFoundResponse(message string, err error, isDebug bool) models.Response {
	return ErrorResponse(http.StatusNotFound, message, err, isDebug)
}

// UnprocessableEntityResponse creates a 422 Unprocessable Entity response
func UnprocessableEntityResponse(message string, err error, isDebug bool) models.Response {
	return ErrorResponse(http.StatusUnprocessableEntity, message, err, isDebug)
}

// InternalServerErrorResponse creates a 500 Internal Server Error response
func InternalServerErrorResponse(message string, err error, isDebug bool) models.Response {
	return ErrorResponse(http.StatusInternalServerError, message, err, isDebug)
}
