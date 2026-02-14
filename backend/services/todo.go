package services

import (
	"context"
	"log/slog"
	"todo-app/dto"
	"todo-app/ent"
	apperrors "todo-app/errors"
	"todo-app/repositories"
)

type TodoServiceFactory func(*slog.Logger) (*TodoService, error)

func ProvideTodoServiceFactory() TodoServiceFactory {
	return func(logger *slog.Logger) (*TodoService, error) {
		repo := repositories.NewTodoRepository(logger)
		return NewTodoService(logger, repo), nil
	}
}

type ITodoRepository interface {
	FetchTodos(ctx context.Context, limit int, offset int, includeDone bool) ([]*ent.Todo, error)
	GetTodoCount(ctx context.Context, includeDone bool) (int, error)
	FindTodo(ctx context.Context, id int) (*ent.Todo, error)
	GetTodoForUpdate(ctx context.Context, id int) (*ent.Todo, error)
	CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error)
	UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error)
	UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error)
	DeleteTodo(ctx context.Context, id int) error
}

func NewTodoService(logger *slog.Logger, repo ITodoRepository) *TodoService {
	return &TodoService{
		logger: logger,
		repo:   repo,
	}
}

type TodoService struct {
	logger *slog.Logger
	repo   ITodoRepository
}

func (s *TodoService) GetTodoSlice(ctx context.Context, currentPage int, limit int, includeDone bool) ([]dto.TodoDto, error) {
	offset := (currentPage - 1) * limit
	todos, err := s.repo.FetchTodos(ctx, limit, offset, includeDone)
	if err != nil {
		return nil, err
	}

	todoDtos := make([]dto.TodoDto, len(todos))
	for i, t := range todos {
		todoDtos[i] = dto.EntityToTodoDto(t)
	}
	return todoDtos, nil
}

func (s *TodoService) CalculatePagination(ctx context.Context, currentPage int, limit int, includeDone bool) (*dto.PaginationDto, error) {
	count, err := s.repo.GetTodoCount(ctx, includeDone)
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

func (s *TodoService) CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error) {
	return s.repo.CreateTodo(ctx, title, description)
}

func (s *TodoService) UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error) {
	client := ent.FromContext(ctx)
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	txCtx := ent.NewTxContext(ctx, tx)

	todo, err := s.repo.GetTodoForUpdate(txCtx, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if todo.DoneAt != nil {
		tx.Rollback()
		// インポートが必要
		return nil, apperrors.ErrTodoAlreadyDone
	}

	if title == nil && description == nil {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return todo, nil
	}

	updatedTodo, err := s.repo.UpdateTodo(txCtx, id, title, description)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updatedTodo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	return s.repo.DeleteTodo(ctx, id)
}

func (s *TodoService) UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error) {
	client := ent.FromContext(ctx)
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	txCtx := ent.NewTxContext(ctx, tx)
	todo, err := s.repo.UpdateDoneStatus(txCtx, id, isDone)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return nil, rerr
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return todo, nil
}
