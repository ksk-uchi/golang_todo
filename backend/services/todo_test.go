package services_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
	"todo-app/ent"
	"todo-app/ent/enttest"
	"todo-app/services"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type spyTodoRepo struct {
	// repositories.TodoRepository (removed embedding)
	fetchTodos       func(limit int, offset int, includeDone bool) ([]*ent.Todo, error)
	count            func(includeDone bool) (int, error)
	create           func(title string, description string) (*ent.Todo, error)
	update           func(id int, title *string, description *string) (*ent.Todo, error)
	updateDoneStatus func(id int, isDone bool) (*ent.Todo, error)
	find             func(id int) (*ent.Todo, error)
	delete           func(id int) error
	getTodoForUpdate func(id int) (*ent.Todo, error)
}

// ... existing methods ...

func (s *spyTodoRepo) FetchTodos(ctx context.Context, limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
	if s.fetchTodos != nil {
		return s.fetchTodos(limit, offset, includeDone)
	}
	return nil, nil
}

func (s *spyTodoRepo) GetTodoCount(ctx context.Context, includeDone bool) (int, error) {
	if s.count != nil {
		return s.count(includeDone)
	}
	return 0, nil
}

func (s *spyTodoRepo) FindTodo(ctx context.Context, id int) (*ent.Todo, error) {
	if s.find != nil {
		return s.find(id)
	}
	return nil, nil
}

func (s *spyTodoRepo) CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error) {
	if s.create != nil {
		return s.create(title, description)
	}
	return nil, nil
}

func (s *spyTodoRepo) UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error) {
	if s.update != nil {
		return s.update(id, title, description)
	}
	return nil, nil
}

func (s *spyTodoRepo) UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error) {
	if s.updateDoneStatus != nil {
		return s.updateDoneStatus(id, isDone)
	}
	return nil, nil
}

func (s *spyTodoRepo) DeleteTodo(ctx context.Context, id int) error {
	if s.delete != nil {
		return s.delete(id)
	}
	return nil
}

func (s *spyTodoRepo) GetTodoForUpdate(ctx context.Context, id int) (*ent.Todo, error) {
	if s.getTodoForUpdate != nil {
		return s.getTodoForUpdate(id)
	}
	return nil, nil
}

// ...

func TestTodoService_UpdateTodo(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	t.Run("リポジトリに更新が正常に反映されること", func(t *testing.T) {
		title := "Updated Title"
		desc := "Updated Description"
		repo := &spyTodoRepo{
			getTodoForUpdate: func(id int) (*ent.Todo, error) {
				return &ent.Todo{ID: id, Title: "Old", Description: "Old", DoneAt: nil}, nil
			},
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
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(ctx, 1, &title, &desc)

		assert.NoError(t, err)
		if err != nil {
			return
		}
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Title", result.Title)
		assert.Equal(t, "Updated Description", result.Description)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		title := "Updated Title"
		repo := &spyTodoRepo{
			getTodoForUpdate: func(id int) (*ent.Todo, error) {
				return &ent.Todo{ID: id, DoneAt: nil}, nil
			},
			update: func(id int, t *string, d *string) (*ent.Todo, error) {
				return nil, errors.New("db error")
			},
		}
		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(ctx, 1, &title, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("更新項目がない場合、更新を行わずID検索を実行すること(UpdateTodoではFindTodoではなくGetTodoForUpdateだけで済むようになったが、ロジック上は)", func(t *testing.T) {
		// New implementation uses GetTodoForUpdate, checks for changes (none), and returns the todo.
		// It does NOT call FindTodo anymore, it reuses the one from GetTodoForUpdate assuming it's fresh enough (it is locked).
		repo := &spyTodoRepo{
			getTodoForUpdate: func(id int) (*ent.Todo, error) {
				return &ent.Todo{
					ID:          id,
					Title:       "Existing Title",
					Description: "Existing Description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
			update: func(id int, t *string, d *string) (*ent.Todo, error) {
				return nil, errors.New("update should not be called")
			},
		}
		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(ctx, 1, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Existing Title", result.Title)
		assert.Equal(t, "Existing Description", result.Description)
	})
}

func TestTodoService_UpdateDoneStatus(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	t.Run("リポジトリに更新が正常に反映されること", func(t *testing.T) {
		repo := &spyTodoRepo{
			updateDoneStatus: func(id int, isDone bool) (*ent.Todo, error) {
				now := time.Now()
				return &ent.Todo{
					ID:          id,
					Title:       "Todo",
					Description: "Desc",
					DoneAt:      &now,
					CreatedAt:   now,
					UpdatedAt:   now,
				}, nil
			},
		}
		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateDoneStatus(ctx, 1, true)

		assert.NoError(t, err)
		if err != nil {
			return
		}
		assert.NotNil(t, result)
		assert.NotNil(t, result.DoneAt)
	})

	t.Run("リポジトリがエラーを返した場合、ロールバックされエラーを返すこと", func(t *testing.T) {
		repo := &spyTodoRepo{
			updateDoneStatus: func(id int, isDone bool) (*ent.Todo, error) {
				return nil, errors.New("db error")
			},
		}
		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateDoneStatus(ctx, 1, true)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
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
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		err := service.DeleteTodo(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := &spyTodoRepo{
			delete: func(id int) error {
				return errors.New("db error")
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		err := service.DeleteTodo(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
