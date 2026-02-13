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
	fetchTodos func(limit int, offset int, includeDone bool) ([]*ent.Todo, error)
	count      func(includeDone bool) (int, error)
	create     func(title string, description string) (*ent.Todo, error)
	update     func(id int, title *string, description *string) (*ent.Todo, error)
	find       func(id int) (*ent.Todo, error)
	delete     func(id int) error
}

func (s *spyTodoRepo) FetchTodos(limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
	if s.fetchTodos != nil {
		return s.fetchTodos(limit, offset, includeDone)
	}
	return nil, nil
}

func (s *spyTodoRepo) GetTodoCount(includeDone bool) (int, error) {
	if s.count != nil {
		return s.count(includeDone)
	}
	return 0, nil
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

func (s *spyTodoRepo) DeleteTodo(id int) error {
	if s.delete != nil {
		return s.delete(id)
	}
	return nil
}

func TestTodoService_GetTodoSlice(t *testing.T) {
	t.Run("リポジトリから取得したデータが正しくDTOに変換されて返ること", func(t *testing.T) {
		repo := &spyTodoRepo{
			fetchTodos: func(limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
				return []*ent.Todo{
					{ID: 2, Title: "Test Task 2", CreatedAt: time.Date(2026, 2, 4, 10, 0, 2, 0, time.Local)},
					{ID: 1, Title: "Test Task 1", CreatedAt: time.Date(2026, 2, 4, 10, 0, 1, 0, time.Local)},
					{ID: 3, Title: "Test Task 3", CreatedAt: time.Date(2026, 2, 4, 10, 0, 0, 0, time.Local)},
				}, nil
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		results, err := service.GetTodoSlice(1, 10, false)

		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "Test Task 2", results[0].Title)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := &spyTodoRepo{
			fetchTodos: func(limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
				return nil, errors.New("db connection error")
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		results, err := service.GetTodoSlice(1, 10, false)

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

func TestTodoService_DeleteTodo(t *testing.T) {
	t.Run("リポジトリの削除が正常に終了すること", func(t *testing.T) {
		repo := &spyTodoRepo{
			delete: func(id int) error {
				return nil
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		err := service.DeleteTodo(1)

		assert.NoError(t, err)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := &spyTodoRepo{
			delete: func(id int) error {
				return errors.New("db error")
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		err := service.DeleteTodo(1)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
