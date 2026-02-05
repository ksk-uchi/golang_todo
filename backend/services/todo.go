package services

import (
	"context"
	"log/slog"
	"todo-app/ent"
	"todo-app/repositories"
)

type TodoServiceFactory func(context.Context, *slog.Logger, *ent.Client) (*TodoService, error)

func ProvideTodoServiceFactory() TodoServiceFactory {
	return func(ctx context.Context, logger *slog.Logger, client *ent.Client) (*TodoService, error) {
		repo := repositories.NewTodoRepository(ctx, client)
		return NewTodoService(ctx, logger, repo), nil
	}
}

type ITodoRepository interface {
	FetchAllTodo() ([]*ent.Todo, error)
	FindTodo(id int) (*ent.Todo, error)
	CreateTodo(title string, description string) (*ent.Todo, error)
	UpdateTodo(id int, title *string, description *string) (*ent.Todo, error)
	DeleteTodo(id int) error
}

func NewTodoService(ctx context.Context, logger *slog.Logger, repo ITodoRepository) *TodoService {
	return &TodoService{
		ctx:    ctx,
		logger: logger,
		repo:   repo,
	}
}

type TodoService struct {
	ctx    context.Context
	logger *slog.Logger
	repo   ITodoRepository
}

func (s *TodoService) GetTodoSlice() ([]*ent.Todo, error) {
	todos, err := s.repo.FetchAllTodo()
	return todos, err
}

func (s *TodoService) CreateTodo(title string, description string) (*ent.Todo, error) {
	return s.repo.CreateTodo(title, description)
}

func (s *TodoService) UpdateTodo(id int, title *string, description *string) (*ent.Todo, error) {
	if title == nil && description == nil {
		return s.repo.FindTodo(id)
	}
	return s.repo.UpdateTodo(id, title, description)
}

func (s *TodoService) DeleteTodo(id int) error {
	return s.repo.DeleteTodo(id)
}
