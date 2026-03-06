package testutils

import (
	"context"
	"todo-app/ent"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTodoFilterHistoryRepository struct {
	mock.Mock
}

func (m *MockTodoFilterHistoryRepository) FetchLatestFilters(ctx context.Context, limit int) ([]*ent.TodoFilterHistory, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ent.TodoFilterHistory), args.Error(1)
}

func (m *MockTodoFilterHistoryRepository) SaveFilterHistory(ctx context.Context, query string, functionName *string, args map[string]interface{}, resultTodoIds []int) (*ent.TodoFilterHistory, error) {
	callArgs := m.Called(ctx, query, functionName, args, resultTodoIds)
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(*ent.TodoFilterHistory), callArgs.Error(1)
}

func (m *MockTodoFilterHistoryRepository) GetFilterHistoryByQueryID(ctx context.Context, queryID uuid.UUID) (*ent.TodoFilterHistory, error) {
	args := m.Called(ctx, queryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ent.TodoFilterHistory), args.Error(1)
}
