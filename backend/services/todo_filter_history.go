package services

import (
	"context"
	"log/slog"
	"todo-app/ent"
	"todo-app/repositories"
)

type ITodoFilterHistoryService interface {
	FetchLatestFilters(ctx context.Context) ([]*ent.TodoFilterHistory, error)
}

type TodoFilterHistoryService struct {
	repo   repositories.ITodoFilterHistoryRepository
	logger *slog.Logger
}

func NewTodoFilterHistoryService(repo repositories.ITodoFilterHistoryRepository, logger *slog.Logger) *TodoFilterHistoryService {
	return &TodoFilterHistoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *TodoFilterHistoryService) FetchLatestFilters(ctx context.Context) ([]*ent.TodoFilterHistory, error) {
	return s.repo.FetchLatestFilters(ctx, 5)
}
