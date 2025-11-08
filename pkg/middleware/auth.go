package middleware

import (
	"github.com/Alfian57/golang-todo/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtUtils *utils.JWTUtils, isDebug bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("Authorization")

		claims, err := jwtUtils.ParseJWT(apiKey)
		if err != nil {
			response := utils.UnauthorizedResponse("Unauthorized", err, isDebug)
			ctx.JSON(response.StatusCode, response)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
