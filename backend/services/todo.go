package services

import (
	"context"
	"todo-app/ent"
	"todo-app/repositories"
)

type TodoServiceFactory func(*ent.Client) *TodoService

func ProvideTodoServiceFactory() TodoServiceFactory {
	return func(client *ent.Client) *TodoService {
		repo := repositories.NewTodoRepository(client)
		return NewTodoService(repo)
	}
}

type ITodoRepository interface {
	FetchAllTodo(ctx context.Context) ([]*ent.Todo, error)
}

type TodoService struct {
	repo ITodoRepository
}

func NewTodoService(repo ITodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) GetTodoSlice(ctx context.Context) ([]*ent.Todo, error) {
	return s.repo.FetchAllTodo(ctx)
}
