package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/Alfian57/golang-todo/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWTUtils provides JWT token operations with injected configuration
type JWTUtils struct {
	appName   string
	jwtSecret []byte
	tokenTTL  time.Duration
}

// NewJWTUtils creates a new JWTUtils instance with the given configuration
func NewJWTUtils(cfg *config.Config) *JWTUtils {
	return &JWTUtils{
		appName:   cfg.App.Name,
		jwtSecret: cfg.JWT.Secret,
		tokenTTL:  cfg.JWT.TTL,
	}
}

// CreateJWT creates a new JWT token for the given user ID
func (j *JWTUtils) CreateJWT(userId string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    j.appName,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

// ParseJWT parses and validates a JWT token string
func (j *JWTUtils) ParseJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return j.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// GetJWTTTL returns the expiration time for a new token
func (j *JWTUtils) GetJWTTTL() time.Time {
	return time.Now().Add(j.tokenTTL)
}

// CreateRefreshToken generates a new refresh token
func CreateRefreshToken() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
