package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrClaimsNotFound    = errors.New("claims not found in context")
	ErrInvalidClaimsType = errors.New("invalid claims type")
	ErrInvalidUserID     = errors.New("invalid user ID in token")
	ErrUserIDNotFound    = errors.New("user ID not found in token")
)

// GetUserIDFromContext extracts and validates user ID from JWT claims in context
// This function should be used in handlers that are protected by AuthMiddleware
func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	// Get claims from context (set by AuthMiddleware)
	claimsInterface, exists := ctx.Get("claims")
	if !exists {
		return uuid.Nil, ErrClaimsNotFound
	}

	// Type assert to *jwt.RegisteredClaims
	claims, ok := claimsInterface.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, ErrInvalidClaimsType
	}

	// Check if Subject (user ID) exists
	if claims.Subject == "" {
		return uuid.Nil, ErrUserIDNotFound
	}

	// Parse userID from Subject field
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrInvalidUserID
	}

	return userID, nil
}

// GetClaimsFromContext extracts JWT claims from context
// Use this if you need access to full claims data (issuer, expiry, etc.)
func GetClaimsFromContext(ctx *gin.Context) (*jwt.RegisteredClaims, error) {
	claimsInterface, exists := ctx.Get("claims")
	if !exists {
		return nil, ErrClaimsNotFound
	}

	claims, ok := claimsInterface.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidClaimsType
	}

	return claims, nil
}
