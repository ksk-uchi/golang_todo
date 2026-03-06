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
	"todo-app/testutils"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ...

func TestTodoService_UpdateTodo(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	t.Run("リポジトリに更新が正常に反映されること", func(t *testing.T) {
		title := "Updated Title"
		desc := "Updated Description"
		repo := new(testutils.MockTodoRepository)
		repo.On("FindTodo", mock.Anything, 1).Return(&ent.Todo{ID: 1, Title: "Old", Description: "Old", DoneAt: nil}, nil)
		repo.On("GetTodoForUpdate", mock.Anything, 1).Return(&ent.Todo{ID: 1, Title: "Old", Description: "Old", DoneAt: nil}, nil)
		repo.On("UpdateTodo", mock.Anything, 1, &title, &desc).Return(&ent.Todo{
			ID:          1,
			Title:       title,
			Description: desc,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)

		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

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
		repo := new(testutils.MockTodoRepository)
		repo.On("FindTodo", mock.Anything, 1).Return(&ent.Todo{ID: 1, DoneAt: nil}, nil)
		repo.On("GetTodoForUpdate", mock.Anything, 1).Return(&ent.Todo{ID: 1, DoneAt: nil}, nil)
		repo.On("UpdateTodo", mock.Anything, 1, &title, (*string)(nil)).Return((*ent.Todo)(nil), errors.New("db error"))

		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(ctx, 1, &title, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("更新項目がない場合、更新を行わずID検索を実行すること(UpdateTodoではFindTodoではなくGetTodoForUpdateだけで済むようになったが、ロジック上は)", func(t *testing.T) {
		repo := new(testutils.MockTodoRepository)
		existingTodo := &ent.Todo{
			ID:          1,
			Title:       "Existing Title",
			Description: "Existing Description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		repo.On("FindTodo", mock.Anything, 1).Return(existingTodo, nil)

		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateTodo(ctx, 1, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Existing Title", result.Title)
		assert.Equal(t, "Existing Description", result.Description)
	})
}

func TestTodoService_UpdateDoneStatus(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	t.Run("リポジトリに更新が正常に反映されること", func(t *testing.T) {
		repo := new(testutils.MockTodoRepository)
		repo.On("GetTodoForUpdate", mock.Anything, 1).Return(&ent.Todo{
			ID:     1,
			DoneAt: nil,
		}, nil)
		now := time.Now()
		repo.On("UpdateDoneStatus", mock.Anything, 1, true).Return(&ent.Todo{
			ID:          1,
			Title:       "Todo",
			Description: "Desc",
			DoneAt:      &now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil)

		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateDoneStatus(ctx, 1, true)

		assert.NoError(t, err)
		if err != nil {
			return
		}
		assert.NotNil(t, result)
		assert.NotNil(t, result.DoneAt)
	})

	t.Run("リポジトリがエラーを返した場合、ロールバックされエラーを返すこと", func(t *testing.T) {
		repo := new(testutils.MockTodoRepository)
		repo.On("GetTodoForUpdate", mock.Anything, 1).Return(&ent.Todo{
			ID:     1,
			DoneAt: nil,
		}, nil)
		repo.On("UpdateDoneStatus", mock.Anything, 1, true).Return((*ent.Todo)(nil), errors.New("db error"))

		ctx := context.Background()
		ctx = ent.NewContext(ctx, client)
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.UpdateDoneStatus(ctx, 1, true)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestTodoService_DeleteTodo(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	t.Run("リポジトリの削除が正常に終了すること", func(t *testing.T) {
		repo := new(testutils.MockTodoRepository)
		repo.On("DeleteTodo", mock.Anything, 1).Return(nil)
		ctx := context.Background()
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		err := service.DeleteTodo(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := new(testutils.MockTodoRepository)
		repo.On("DeleteTodo", mock.Anything, 1).Return(errors.New("db error"))
		ctx := context.Background()
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		err := service.DeleteTodo(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestTodoService_FetchTodosByIds(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	t.Run("リポジトリの FetchTodosByIds を呼び出し、結果を返すこと", func(t *testing.T) {
		repo := new(testutils.MockTodoRepository)
		ids := []int{1, 2, 3}
		expectedTodos := []*ent.Todo{{ID: 1}, {ID: 2}, {ID: 3}}
		repo.On("FetchTodosByIds", mock.Anything, ids).Return(expectedTodos, nil)

		ctx := context.Background()
		service := services.NewTodoService(client, slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		result, err := service.FetchTodosByIds(ctx, ids)

		assert.NoError(t, err)
		assert.Equal(t, expectedTodos, result)
		repo.AssertExpectations(t)
	})
}
