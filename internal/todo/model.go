package todo

import (
	"github.com/Alfian57/golang-todo/common/models"
	"github.com/google/uuid"
)

type Todo struct {
	models.Base
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      uuid.UUID `json:"user_id"`
}
