package testutils

import (
	"context"
	"todo-app/ent"

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
