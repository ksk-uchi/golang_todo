package services

import (
	"context"
	"log/slog"
	"todo-app/dto"
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
	FetchTodos(limit int, offset int) ([]*ent.Todo, error)
	GetTodoCount() (int, error)
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

func (s *TodoService) GetTodoSlice(currentPage int, limit int) ([]*ent.Todo, error) {
	offset := (currentPage - 1) * limit
	todos, err := s.repo.FetchTodos(limit, offset)
	return todos, err
}

func (s *TodoService) CalculatePagination(currentPage int, limit int) (*dto.PaginationDto, error) {
	count, err := s.repo.GetTodoCount()
	if err != nil {
		return nil, err
	}

	totalPages := (count + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.PaginationDto{
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
		Limit:       limit,
	}, nil
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
