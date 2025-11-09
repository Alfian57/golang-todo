package auth

import (
	"time"

	"github.com/Alfian57/golang-todo/common/models"
	"github.com/google/uuid"
)

type User struct {
	models.Base
	Username string
	Password string
}

type RefreshToken struct {
	models.Base
	Token     string `gorm:"index"`
	UserID    uuid.UUID
	ExpiredAt time.Time
}
