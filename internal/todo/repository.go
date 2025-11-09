package todo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return TodoRepository{
		db: db,
	}
}

func (repository TodoRepository) FindAllTodoByUserID(ctx context.Context, userID string) ([]Todo, error) {
	var todos []Todo
	err := repository.db.Where("user_id = ?", userID).Find(&todos).Error
	return todos, err
}

func (repository TodoRepository) CreateTodo(ctx context.Context, todo *Todo) error {
	return repository.db.Create(todo).Error
}

func (repository TodoRepository) FindTodoByID(ctx context.Context, id uint) (*Todo, error) {
	todo, err := gorm.G[Todo](repository.db).First(ctx)
	return &todo, err
}

func (repository TodoRepository) FindTodoByIDAndUserID(ctx context.Context, todoID uuid.UUID, userID uuid.UUID) (*Todo, error) {
	var todo Todo
	err := repository.db.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error
	return &todo, err
}

func (repository TodoRepository) UpdateTodo(ctx context.Context, todo *Todo) error {
	return repository.db.Save(todo).Error
}

func (repository TodoRepository) DeleteTodo(ctx context.Context, todo *Todo) error {
	return repository.db.Delete(todo).Error
}
