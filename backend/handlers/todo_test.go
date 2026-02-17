package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
	"todo-app/di"
	"todo-app/dto"
	"todo-app/ent"
	"todo-app/ent/todo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v5"
	_ "github.com/mattn/go-sqlite3" // テスト実行にドライバが必要
	"github.com/stretchr/testify/assert"
)

func TestTodoHandler_ListTodo_Integration(t *testing.T) {
	t.Run("Todo 一覧取得", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
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
					SetUpdatedAt(d.CreatedAt).
					SetUser(user).
					SaveX(context.Background())
			}

			req, rec := createAuthenticatedRequest(t, http.MethodGet, "/todo", "", user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			var res dto.ListTodoResponseDto
			_ = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Len(t, res.Data, 3)
			assert.Equal(t, "Test Title 2", res.Data[0].Title)
			assert.Equal(t, "Test Title 1", res.Data[1].Title)
			assert.Equal(t, "Test Title 3", res.Data[2].Title)

			// Check Pagination
			assert.NotNil(t, res.Pagination)
			assert.Equal(t, 1, res.Pagination.TotalPages)
			assert.Equal(t, 1, res.Pagination.CurrentPage)
			assert.False(t, res.Pagination.HasNext)
			assert.False(t, res.Pagination.HasPrev)
			assert.Equal(t, 20, res.Pagination.Limit)
		})
	})
}

func TestTodoHandler_ListTodo_IncludeDone_Integration(t *testing.T) {
	withTx(t, func(client *ent.Client) {
		e := echo.New()
		app, err := di.InitializeTestApp(e, client)
		assert.NoError(t, err)

		app.Router.Setup(e)

		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

		// 1. Active Todo
		client.Todo.Create().
			SetTitle("Active Todo 1").
			SetDescription("Desc").
			SetUser(user).
			SaveX(context.Background())

		// 2. Done Todo
		client.Todo.Create().
			SetTitle("Done Todo").
			SetDescription("Desc").
			SetDoneAt(time.Now()).
			SetUser(user).
			SaveX(context.Background())

		// 3. Active Todo
		client.Todo.Create().
			SetTitle("Active Todo 2").
			SetDescription("Desc").
			SetUser(user).
			SaveX(context.Background())

		t.Run("Default returns only active todos", func(t *testing.T) {
			req, rec := createAuthenticatedRequest(t, http.MethodGet, "/todo", "", user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			var res dto.ListTodoResponseDto
			_ = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Len(t, res.Data, 2)
			for _, todo := range res.Data {
				assert.Contains(t, todo.Title, "Active")
			}
		})

		t.Run("include_done=true returns all todos", func(t *testing.T) {
			req, rec := createAuthenticatedRequest(t, http.MethodGet, "/todo?include_done=true", "", user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			var res dto.ListTodoResponseDto
			_ = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Len(t, res.Data, 3)
		})

		t.Run("include_done=false returns only active todos", func(t *testing.T) {
			req, rec := createAuthenticatedRequest(t, http.MethodGet, "/todo?include_done=false", "", user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			var res dto.ListTodoResponseDto
			_ = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Len(t, res.Data, 2)
		})
	})
}

func TestTodoHandler_CreateTodo_Integration(t *testing.T) {
	t.Run("Todo 新規作成", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)

			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
			body := `{"title": "New Todo", "description": "New Description"}`

			req, rec := createAuthenticatedRequest(t, http.MethodPost, "/todo", body, user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusCreated, rec.Code)

			var res dto.TodoDto
			_ = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Equal(t, "New Todo", res.Title)
			assert.Equal(t, "New Description", res.Description)

			// DBにも正しく保存されていることを確認
			count, err := client.Todo.Query().Count(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, 1, count)
		})
	})

	t.Run("Todo 新規作成 バリデーションエラー", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)
			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

			body := `{"title": "", "description": "New Description"}`

			req, rec := createAuthenticatedRequest(t, http.MethodPost, "/todo", body, user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			expected := `{"error":{"title":"タイトルは必須です"}}`
			assert.JSONEq(t, expected, rec.Body.String())
		})
	})

	t.Run("Todo 新規作成 未認証", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)

			body := `{"title": "New Todo", "description": "New Description"}`

			req, rec := createAuthenticatedRequest(t, http.MethodPost, "/todo", body, 0)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.JSONEq(t, `{"error":"missing token"}`, rec.Body.String())
		})
	})
}

func TestTodoHandler_UpdateTodo_Integration(t *testing.T) {
	if os.Getenv("TEST_WITH_REAL_DB") == "" {
		t.Skip("TEST_WITH_REAL_DB is not set. Skipping integration test that requires real DB (e.g. MySQL) for SELECT FOR UPDATE.")
	}

	t.Run("Todo 更新", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
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

			req, rec := createAuthenticatedRequest(t, http.MethodPatch, fmt.Sprintf("/todo/%d", todo.ID), body, user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			var res dto.TodoDto
			_ = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Equal(t, "Updated Title", res.Title)
			assert.Equal(t, "Updated Description", res.Description)
			assert.Equal(t, todo.ID, res.ID)

			// Check DB
			updatedTodo := client.Todo.GetX(context.Background(), todo.ID)
			assert.Equal(t, "Updated Title", updatedTodo.Title)
			assert.Equal(t, "Updated Description", updatedTodo.Description)
		})
	})

	t.Run("Todo 更新 存在しないID", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)
			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

			body := `{"title": "Updated Title"}`

			req, rec := createAuthenticatedRequest(t, http.MethodPatch, "/todo/999", body, user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code)
		})
	})

	t.Run("Todo 更新 他人のTodo", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)

			user1 := client.User.Create().SetName("user1").SetEmail("user1").SetPassword("user1").SaveX(context.Background())
			user2 := client.User.Create().SetName("user2").SetEmail("user2").SetPassword("user2").SaveX(context.Background())

			todo := client.Todo.Create().
				SetTitle("User1 Todo").
				SetDescription("User1 Description").
				SetUser(user1).
				SaveX(context.Background())

			body := `{"title": "Updated Title"}`

			req, rec := createAuthenticatedRequest(t, http.MethodPatch, fmt.Sprintf("/todo/%d", todo.ID), body, user2.ID)
			e.ServeHTTP(rec, req)

			// Expect Not Found because repository filters by user
			assert.Equal(t, http.StatusNotFound, rec.Code)
		})
	})
	t.Run("Todo 更新 完了済み", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)

			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
			todo := client.Todo.Create().
				SetTitle("Done Todo").
				SetDescription("Description").
				SetDoneAt(time.Now()).
				SetUser(user).
				SaveX(context.Background())

			body := `{"title": "Updated Title", "description": "Updated Description"}`

			req, rec := createAuthenticatedRequest(t, http.MethodPatch, fmt.Sprintf("/todo/%d", todo.ID), body, user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "cannot update a completed todo")
		})
	})
}

func TestTodoHandler_DeleteTodo_Integration(t *testing.T) {
	t.Run("Todo 削除", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
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

			req, rec := createAuthenticatedRequest(t, http.MethodDelete, fmt.Sprintf("/todo/%d", todo.ID), "", user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNoContent, rec.Code)

			// DBから消えているか確認
			count, err := client.Todo.Query().Count(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	})

	t.Run("Todo 削除 存在しないID", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)
			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

			req, rec := createAuthenticatedRequest(t, http.MethodDelete, "/todo/999", "", user.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNoContent, rec.Code)
		})
	})

	t.Run("Todo 削除 他人のTodo", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)

			user1 := client.User.Create().SetName("user1").SetEmail("user1").SetPassword("user1").SaveX(context.Background())
			user2 := client.User.Create().SetName("user2").SetEmail("user2").SetPassword("user2").SaveX(context.Background())

			targetTodo := client.Todo.Create().
				SetTitle("To Be Deleted").
				SetDescription("Description").
				SetUser(user1).
				SaveX(context.Background())

			req, rec := createAuthenticatedRequest(t, http.MethodDelete, fmt.Sprintf("/todo/%d", targetTodo.ID), "", user2.ID)
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNoContent, rec.Code)

			// Ensure it was NOT deleted
			exists := client.Todo.Query().Where(todo.ID(targetTodo.ID)).ExistX(context.Background())
			assert.True(t, exists)
		})
	})
}

func TestTodoHandler_UpdateDoneStatus(t *testing.T) {
	if os.Getenv("TEST_WITH_REAL_DB") == "" {
		t.Skip("TEST_WITH_REAL_DB is not set. Skipping integration test that requires real DB (e.g. MySQL) for SELECT FOR UPDATE.")
	}

	t.Run("成功時", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)

			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
			todo := client.Todo.Create().
				SetTitle("Todo").
				SetDescription("Desc").
				SetUser(user).
				SaveX(context.Background())

			reqBody := map[string]interface{}{"is_done": true}
			body, _ := json.Marshal(reqBody)
			req, rec := createAuthenticatedRequest(t, http.MethodPut, fmt.Sprintf("/todo/%d/done", todo.ID), string(body), user.ID)

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "\"done_at\"")

			// Verify DB
			updatedTodo := client.Todo.GetX(context.Background(), todo.ID)
			assert.NotNil(t, updatedTodo.DoneAt)
		})
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)
			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
			todo := client.Todo.Create().
				SetTitle("Todo").
				SetDescription("Desc").
				SetUser(user).
				SaveX(context.Background())

			reqBody := map[string]interface{}{} // is_done missing
			body, _ := json.Marshal(reqBody)
			req, rec := createAuthenticatedRequest(t, http.MethodPut, fmt.Sprintf("/todo/%d/done", todo.ID), string(body), user.ID)

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("Todoが見つからない場合", func(t *testing.T) {
		withTx(t, func(client *ent.Client) {
			e := echo.New()
			app, err := di.InitializeTestApp(e, client)
			assert.NoError(t, err)

			app.Router.Setup(e)
			user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

			reqBody := map[string]interface{}{"is_done": true}
			body, _ := json.Marshal(reqBody)
			req, rec := createAuthenticatedRequest(t, http.MethodPut, "/todo/999/done", string(body), user.ID)

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code)
		})
	})
}
