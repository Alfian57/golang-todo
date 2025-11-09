package todo

import (
	"github.com/google/uuid"
)

// General TodoResponse
type TodoResponse struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   string    `json:"created_at"`
}

// List Todos
type GetTodosResponse struct {
	Todos []TodoResponse `json:"todos"`
}

// Create Todo
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,max=255,min=1"`
	Description string `json:"description" validate:"max=1000"`
}
type CreateTodoResponse struct {
	Todo TodoResponse `json:"todo"`
}

// Update Todo
type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required,max=255,min=1"`
	Description string `json:"description" validate:"max=1000"`
	Completed   bool   `json:"completed"`
}
type UpdateTodoResponse struct {
	Todo TodoResponse `json:"todo"`
}
