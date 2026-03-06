package services

import (
	"context"
	"log/slog"
	"todo-app/ent"
	"todo-app/repositories"

	"github.com/google/uuid"
)

type ITodoFilterHistoryService interface {
	FetchLatestFilters(ctx context.Context) ([]*ent.TodoFilterHistory, error)
	SaveFilterHistory(ctx context.Context, query string, functionName *string, args map[string]interface{}, resultTodoIds []int) (*ent.TodoFilterHistory, error)
	GetFilterHistoryByQueryID(ctx context.Context, queryID uuid.UUID) (*ent.TodoFilterHistory, error)
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

func (s *TodoFilterHistoryService) SaveFilterHistory(ctx context.Context, query string, functionName *string, args map[string]interface{}, resultTodoIds []int) (*ent.TodoFilterHistory, error) {
	return s.repo.SaveFilterHistory(ctx, query, functionName, args, resultTodoIds)
}

func (s *TodoFilterHistoryService) GetFilterHistoryByQueryID(ctx context.Context, queryID uuid.UUID) (*ent.TodoFilterHistory, error) {
	return s.repo.GetFilterHistoryByQueryID(ctx, queryID)
}
