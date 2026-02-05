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
	create   func(title string, description string) (*ent.Todo, error)
	update   func(id int, title *string, description *string) (*ent.Todo, error)
	find     func(id int) (*ent.Todo, error)
}

func (s *spyTodoRepo) FetchAllTodo() ([]*ent.Todo, error) {
	if s.fetchAll != nil {
		return s.fetchAll()
	}
	return nil, nil
}

func (s *spyTodoRepo) FindTodo(id int) (*ent.Todo, error) {
	if s.find != nil {
		return s.find(id)
	}
	return nil, nil
}

func (s *spyTodoRepo) CreateTodo(title string, description string) (*ent.Todo, error) {
	if s.create != nil {
		return s.create(title, description)
	}
	return nil, nil
}

func (s *spyTodoRepo) UpdateTodo(id int, title *string, description *string) (*ent.Todo, error) {
	if s.update != nil {
		return s.update(id, title, description)
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

func TestTodoService_CreateTodo(t *testing.T) {
	t.Run("リポジトリにデータが正常に保存されること", func(t *testing.T) {
		repo := &spyTodoRepo{
			create: func(title string, description string) (*ent.Todo, error) {
				return &ent.Todo{
					ID:          1,
					Title:       title,
					Description: description,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.CreateTodo("New Task", "New Description")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Task", result.Title)
		assert.Equal(t, "New Description", result.Description)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := &spyTodoRepo{
			create: func(title string, description string) (*ent.Todo, error) {
				return nil, errors.New("db connection error")
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.CreateTodo("New Task", "New Description")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db connection error", err.Error())
	})
}

func TestTodoService_UpdateTodo(t *testing.T) {
	t.Run("リポジトリに更新が正常に反映されること", func(t *testing.T) {
		title := "Updated Title"
		desc := "Updated Description"
		repo := &spyTodoRepo{
			update: func(id int, t *string, d *string) (*ent.Todo, error) {
				return &ent.Todo{
					ID:          1,
					Title:       *t,
					Description: *d,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(1, &title, &desc)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Title", result.Title)
		assert.Equal(t, "Updated Description", result.Description)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		title := "Updated Title"
		repo := &spyTodoRepo{
			update: func(id int, t *string, d *string) (*ent.Todo, error) {
				return nil, errors.New("db error")
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(1, &title, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("更新項目がない場合、更新を行わずID検索を実行すること", func(t *testing.T) {
		repo := &spyTodoRepo{
			update: func(id int, t *string, d *string) (*ent.Todo, error) {
				return nil, errors.New("update should not be called")
			},
			find: func(id int) (*ent.Todo, error) {
				return &ent.Todo{
					ID:          id,
					Title:       "Existing Title",
					Description: "Existing Description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(1, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Existing Title", result.Title)
		assert.Equal(t, "Existing Description", result.Description)
	})
}
