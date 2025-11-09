package todo

import (
	"context"
	"time"

	"github.com/Alfian57/golang-todo/common/models"
	"github.com/Alfian57/golang-todo/pkg/logger"
	"github.com/Alfian57/golang-todo/pkg/utils"
	"github.com/google/uuid"
)

type TodoService struct {
	todoRepository TodoRepository
	log            logger.Logger
	isDebug        bool
}

func NewTodoService(todoRepository TodoRepository, log logger.Logger, isDebug bool) TodoService {
	return TodoService{
		todoRepository: todoRepository,
		log:            log,
		isDebug:        isDebug,
	}
}

func (service TodoService) GetAll(ctx context.Context, userID uuid.UUID) models.Response {
	todos, err := service.todoRepository.FindAllTodoByUserID(ctx, userID.String())
	if err != nil {
		service.log.Error("Failed to get todos",
			logger.F("operation", "Get all todos"),
			logger.F("user_id", userID.String()),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to get todos", err, service.isDebug)
	}

	todosResponse := make([]TodoResponse, 0, len(todos))
	for _, todo := range todos {
		todoResponse := TodoResponse{
			Id:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		todosResponse = append(todosResponse, todoResponse)
	}
	responseData := GetTodosResponse{
		Todos: todosResponse,
	}
	return utils.OkResponse("Todos retrieved successfully", responseData)
}

func (service TodoService) Create(ctx context.Context, title string, description string, userID uuid.UUID) models.Response {
	todo := Todo{
		Title:       title,
		Description: description,
		Completed:   false,
		UserID:      userID,
	}
	err := service.todoRepository.CreateTodo(ctx, &todo)
	if err != nil {
		service.log.Error("Failed to create todo",
			logger.F("operation", "Create todo"),
			logger.F("user_id", userID),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to create todo", err, service.isDebug)
	}

	responseData := CreateTodoResponse{
		Todo: TodoResponse{
			Id:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		},
	}
	return utils.CreatedResponse("Success to create todo", responseData)
}

func (service TodoService) Update(ctx context.Context, todoID uuid.UUID, title string, description string, completed bool, userID uuid.UUID) models.Response {
	todo, err := service.todoRepository.FindTodoByIDAndUserID(ctx, todoID, userID)
	if err != nil {
		service.log.Error("Failed to find todo",
			logger.F("operation", "Update todo"),
			logger.F("todo_id", todoID.String()),
			logger.F("user_id", userID.String()),
			logger.F("error", err),
		)
		return utils.NotFoundResponse("Todo not found", err, service.isDebug)
	}

	todo.Title = title
	todo.Description = description
	todo.Completed = completed

	err = service.todoRepository.UpdateTodo(ctx, todo)
	if err != nil {
		service.log.Error("Failed to update todo",
			logger.F("operation", "Update todo"),
			logger.F("todo_id", todoID.String()),
			logger.F("user_id", userID.String()),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to update todo", err, service.isDebug)
	}

	responseData := UpdateTodoResponse{
		Todo: TodoResponse{
			Id:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}
	return utils.OkResponse("Todo updated successfully", responseData)
}

func (service TodoService) Delete(ctx context.Context, todoID uuid.UUID, userID uuid.UUID) models.Response {
	todo, err := service.todoRepository.FindTodoByIDAndUserID(ctx, todoID, userID)
	if err != nil {
		service.log.Error("Failed to find todo",
			logger.F("operation", "Delete todo"),
			logger.F("todo_id", todoID.String()),
			logger.F("user_id", userID.String()),
			logger.F("error", err),
		)
		return utils.NotFoundResponse("Todo not found", err, service.isDebug)
	}

	err = service.todoRepository.DeleteTodo(ctx, todo)
	if err != nil {
		service.log.Error("Failed to delete todo",
			logger.F("operation", "Delete todo"),
			logger.F("todo_id", todoID.String()),
			logger.F("user_id", userID.String()),
			logger.F("error", err),
		)
		return utils.InternalServerErrorResponse("Failed to delete todo", err, service.isDebug)
	}

	return utils.OkResponse("Todo deleted successfully", nil)
}
