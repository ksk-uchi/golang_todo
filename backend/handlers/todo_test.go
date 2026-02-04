package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
