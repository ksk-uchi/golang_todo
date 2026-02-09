package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"todo-app/di"
	"todo-app/dto"
	"todo-app/ent/enttest"

	"github.com/labstack/echo/v5"
	_ "github.com/mattn/go-sqlite3" // テスト実行にドライバが必要
	"github.com/stretchr/testify/assert"
)

func TestTodoHandler_ListTodo_Integration(t *testing.T) {
	t.Run("Todo 一覧取得", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		type dummyData struct {
			Title       string
			Description string
			CreatedAt   time.Time
		}
		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
		data := []dummyData{
			{
				Title:       "Test Title 1",
				Description: "Test Description 1",
				CreatedAt:   time.Date(2026, 2, 4, 10, 0, 1, 0, time.Local),
			},
			{
				Title:       "Test Title 2",
				Description: "Test Description 2",
				CreatedAt:   time.Date(2026, 2, 4, 10, 0, 2, 0, time.Local),
			},
			{
				Title:       "Test Title 3",
				Description: "Test Description 3",
				CreatedAt:   time.Date(2026, 2, 4, 10, 0, 0, 0, time.Local),
			},
		}

		for _, d := range data {
			client.Todo.Create().
				SetTitle(d.Title).
				SetDescription(d.Description).
				SetCreatedAt(d.CreatedAt).
				SetUser(user).
				SaveX(context.Background())
		}

		req := httptest.NewRequest(http.MethodGet, "/todo", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res []dto.TodoDto
		json.Unmarshal(rec.Body.Bytes(), &res)
		assert.Len(t, res, 3)
		assert.Equal(t, "Test Title 2", res[0].Title)
		assert.Equal(t, "Test Title 1", res[1].Title)
		assert.Equal(t, "Test Title 3", res[2].Title)
	})
}

func TestTodoHandler_CreateTodo_Integration(t *testing.T) {
	t.Run("Todo 新規作成", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		// TODO: ログイン機能実装後に修正する
		client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
		body := `{"title": "New Todo", "description": "New Description"}`
		req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var res dto.TodoDto
		json.Unmarshal(rec.Body.Bytes(), &res)
		assert.Equal(t, "New Todo", res.Title)
		assert.Equal(t, "New Description", res.Description)

		// DBにも正しく保存されていることを確認
		count, err := client.Todo.Query().Count(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Todo 新規作成 バリデーションエラー", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		body := `{"title": "", "description": "New Description"}`
		req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expected := `{"error":{"title":"タイトルは必須です"}}`
		assert.JSONEq(t, expected, rec.Body.String())
	})
}

func TestTodoHandler_UpdateTodo_Integration(t *testing.T) {
	t.Run("Todo 更新", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
		todo := client.Todo.Create().
			SetTitle("Old Title").
			SetDescription("Old Description").
			SetUser(user).
			SaveX(context.Background())

		body := `{"title": "Updated Title", "description": "Updated Description"}`
		req := httptest.NewRequest(http.MethodPatch, "/todo/1", strings.NewReader(body)) // Assuming ID 1

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res dto.TodoDto
		json.Unmarshal(rec.Body.Bytes(), &res)
		assert.Equal(t, "Updated Title", res.Title)
		assert.Equal(t, "Updated Description", res.Description)
		assert.Equal(t, todo.ID, res.ID)

		// Check DB
		updatedTodo := client.Todo.GetX(context.Background(), todo.ID)
		assert.Equal(t, "Updated Title", updatedTodo.Title)
		assert.Equal(t, "Updated Description", updatedTodo.Description)
	})

	t.Run("Todo 更新 存在しないID", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		body := `{"title": "Updated Title"}`
		req := httptest.NewRequest(http.MethodPatch, "/todo/999", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTodoHandler_DeleteTodo_Integration(t *testing.T) {
	t.Run("Todo 削除", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
		todo := client.Todo.Create().
			SetTitle("To Be Deleted").
			SetDescription("Description").
			SetUser(user).
			SaveX(context.Background())

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/todo/%d", todo.ID), nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// DBから消えているか確認
		count, err := client.Todo.Query().Count(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Todo 削除 存在しないID", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		req := httptest.NewRequest(http.MethodDelete, "/todo/999", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
