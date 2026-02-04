package services_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
	"todo-app/ent"
	"todo-app/repositories"
	"todo-app/services"

	"github.com/stretchr/testify/assert"
)

type spyTodoRepo struct {
	repositories.TodoRepository
	fetchAll func() ([]*ent.Todo, error)
}

func (s *spyTodoRepo) FetchAllTodo() ([]*ent.Todo, error) {
	if s.fetchAll != nil {
		return s.fetchAll()
	}
	return nil, nil
}

func TestTodoService_GetTodoSlice(t *testing.T) {
	t.Run("リポジトリから取得したデータがそのまま返ること", func(t *testing.T) {
		repo := &spyTodoRepo{
			fetchAll: func() ([]*ent.Todo, error) {
				return []*ent.Todo{
					{ID: 2, Title: "Test Task 2", CreatedAt: time.Date(2026, 2, 4, 10, 0, 2, 0, time.Local)},
					{ID: 1, Title: "Test Task 1", CreatedAt: time.Date(2026, 2, 4, 10, 0, 1, 0, time.Local)},
					{ID: 3, Title: "Test Task 3", CreatedAt: time.Date(2026, 2, 4, 10, 0, 0, 0, time.Local)},
				}, nil
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		results, err := service.GetTodoSlice()

		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "Test Task 2", results[0].Title)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := &spyTodoRepo{
			fetchAll: func() ([]*ent.Todo, error) {
				return nil, errors.New("db connection error")
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		results, err := service.GetTodoSlice()

		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Equal(t, "db connection error", err.Error())
	})
}
