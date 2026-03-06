package testutils

import (
	"context"
	"time"
	"todo-app/ent"

	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) FetchTodos(ctx context.Context, limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
	args := m.Called(ctx, limit, offset, includeDone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ent.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetTodoCount(ctx context.Context, includeDone bool) (int, error) {
	args := m.Called(ctx, includeDone)
	return args.Int(0), args.Error(1)
}

func (m *MockTodoRepository) FindTodo(ctx context.Context, id int) (*ent.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ent.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetTodoForUpdate(ctx context.Context, id int) (*ent.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ent.Todo), args.Error(1)
}

func (m *MockTodoRepository) CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error) {
	args := m.Called(ctx, title, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ent.Todo), args.Error(1)
}

func (m *MockTodoRepository) UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error) {
	args := m.Called(ctx, id, title, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ent.Todo), args.Error(1)
}

func (m *MockTodoRepository) UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error) {
	args := m.Called(ctx, id, isDone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ent.Todo), args.Error(1)
}

func (m *MockTodoRepository) DeleteTodo(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTodoRepository) FetchTodosByDoneAt(ctx context.Context, doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error) {
	args := m.Called(ctx, doneFrom, doneTo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ent.Todo), args.Error(1)
}

func (m *MockTodoRepository) FetchTodosByIds(ctx context.Context, ids []int) ([]*ent.Todo, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ent.Todo), args.Error(1)
}
