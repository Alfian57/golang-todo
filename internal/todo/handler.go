package todo

import (
	"github.com/Alfian57/golang-todo/pkg/logger"
	"github.com/Alfian57/golang-todo/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	todoService TodoService
	log         logger.Logger
}

func NewTodoHandler(todoService TodoService, log logger.Logger) TodoHandler {
	return TodoHandler{
		todoService: todoService,
		log:         log,
	}
}

// @Summary      Get all todos
// @Description  Get all todos for the authenticated user
// @Tags         Todo
// @Produce      json
// @Security	 BearerAuth
// @Success      200  {object}  models.Response{data=todo.GetTodosResponse}
// @Failure      401  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /todo [get]
func (handler TodoHandler) GetAll(ctx *gin.Context) {
	// Get user ID from JWT claims
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		handler.log.Warn("Failed to get user ID from context",
			logger.F("operation", "Get all todos"),
			logger.F("error", err),
		)
		response := utils.UnauthorizedResponse("Unauthorized", err, false)
		ctx.JSON(response.StatusCode, response)
		return
	}

	response := handler.todoService.GetAll(ctx, userID)
	if response.StatusCode != 200 {
		handler.log.Warn("Get all todos request failed",
			logger.F("operation", "Get all todos"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("user_id", userID.String()),
		)
	}

	ctx.JSON(response.StatusCode, response)
}

// @Summary      Create todo
// @Description  Create a new todo for the authenticated user
// @Tags         Todo
// @Accept       json
// @Produce      json
// @Param        body  body      CreateTodoRequest  true  "Create Todo Request"
// @Security	 BearerAuth
// @Success      201  {object}  models.Response{data=todo.CreateTodoResponse}
// @Failure      401  {object}  models.Response
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /todo [post]
func (handler TodoHandler) Create(ctx *gin.Context) {
	var req CreateTodoRequest
	if !utils.ValidateRequest(ctx, &req) {
		return
	}

	// Get user ID from JWT claims
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		handler.log.Warn("Failed to get user ID from context",
			logger.F("operation", "Create todo"),
			logger.F("error", err),
		)
		response := utils.UnauthorizedResponse("Unauthorized", err, false)
		ctx.JSON(response.StatusCode, response)
		return
	}

	response := handler.todoService.Create(ctx, req.Title, req.Description, userID)
	if response.StatusCode != 201 {
		handler.log.Warn("Create todo request failed",
			logger.F("operation", "Create todo"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("user_id", userID.String()),
		)
	}

	ctx.JSON(response.StatusCode, response)
}

// @Summary      Update todo
// @Description  Update an existing todo for the authenticated user
// @Tags         Todo
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "Todo ID"
// @Param        body  body      UpdateTodoRequest  true  "Update Todo Request"
// @Security	 BearerAuth
// @Success      200  {object}  models.Response{data=todo.UpdateTodoResponse}
// @Failure      401  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /todo/{id} [put]
func (handler TodoHandler) Update(ctx *gin.Context) {
	var req UpdateTodoRequest
	if !utils.ValidateRequest(ctx, &req) {
		return
	}

	// Get todo ID from URL parameter
	todoIDStr := ctx.Param("id")
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		handler.log.Warn("Invalid todo ID",
			logger.F("operation", "Update todo"),
			logger.F("todo_id", todoIDStr),
			logger.F("error", err),
		)
		response := utils.UnprocessableEntityResponse("Invalid todo ID", err, false)
		ctx.JSON(response.StatusCode, response)
		return
	}

	// Get user ID from JWT claims
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		handler.log.Warn("Failed to get user ID from context",
			logger.F("operation", "Update todo"),
			logger.F("error", err),
		)
		response := utils.UnauthorizedResponse("Unauthorized", err, false)
		ctx.JSON(response.StatusCode, response)
		return
	}

	response := handler.todoService.Update(ctx, todoID, req.Title, req.Description, req.Completed, userID)
	if response.StatusCode != 200 {
		handler.log.Warn("Update todo request failed",
			logger.F("operation", "Update todo"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("todo_id", todoID.String()),
			logger.F("user_id", userID.String()),
		)
	}

	ctx.JSON(response.StatusCode, response)
}

// @Summary      Delete todo
// @Description  Delete an existing todo for the authenticated user
// @Tags         Todo
// @Produce      json
// @Param        id   path      string  true  "Todo ID"
// @Security	 BearerAuth
// @Success      200  {object}  models.Response
// @Failure      401  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      422  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /todo/{id} [delete]
func (handler TodoHandler) Delete(ctx *gin.Context) {
	// Get todo ID from URL parameter
	todoIDStr := ctx.Param("id")
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		handler.log.Warn("Invalid todo ID",
			logger.F("operation", "Delete todo"),
			logger.F("todo_id", todoIDStr),
			logger.F("error", err),
		)
		response := utils.UnprocessableEntityResponse("Invalid todo ID", err, false)
		ctx.JSON(response.StatusCode, response)
		return
	}

	// Get user ID from JWT claims
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		handler.log.Warn("Failed to get user ID from context",
			logger.F("operation", "Delete todo"),
			logger.F("error", err),
		)
		response := utils.UnauthorizedResponse("Unauthorized", err, false)
		ctx.JSON(response.StatusCode, response)
		return
	}

	response := handler.todoService.Delete(ctx, todoID, userID)
	if response.StatusCode != 200 {
		handler.log.Warn("Delete todo request failed",
			logger.F("operation", "Delete todo"),
			logger.F("status_code", response.StatusCode),
			logger.F("message", response.Message),
			logger.F("todo_id", todoID.String()),
			logger.F("user_id", userID.String()),
		)
	}

	ctx.JSON(response.StatusCode, response)
}
