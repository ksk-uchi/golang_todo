package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"todo-app/di"
	"todo-app/dto"
	"todo-app/ent"
	"todo-app/ent/enttest"
	"todo-app/ent/todo"
	"todo-app/handlers"
	"todo-app/repositories"
	"todo-app/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	_ "github.com/mattn/go-sqlite3" // テスト実行にドライバが必要
	"github.com/stretchr/testify/assert"
)

type spyTodoRepo struct {
	repositories.TodoRepository
	fetchTodos       func(limit int, offset int, includeDone bool) ([]*ent.Todo, error)
	count            func(includeDone bool) (int, error)
	create           func(title string, description string) (*ent.Todo, error)
	update           func(id int, title *string, description *string) (*ent.Todo, error)
	updateDoneStatus func(id int, isDone bool) (*ent.Todo, error)
	find             func(id int) (*ent.Todo, error)
	delete           func(id int) error
}

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

func TestMain(m *testing.M) {
	if err := godotenv.Load("../envs/test.env"); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

type avoidCSRF struct{}

func (a avoidCSRF) SetCSRFCookie(req *http.Request) {
	cookie := &http.Cookie{Name: "csrf_token", Value: "test"}
	req.AddCookie(cookie)
}

func (a avoidCSRF) SetCSRFHeader(req *http.Request) {
	req.Header.Set("X-CSRF-Token", "test")
}

func createToken(t *testing.T, userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID), // jwt deserializes numbers as float64
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	assert.NoError(t, err)
	return tokenString
}

func createAuthenticatedRequest(t *testing.T, method, target, body string, userID int) (*http.Request, *httptest.ResponseRecorder) {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	if userID != 0 {
		cookie := &http.Cookie{Name: "token", Value: createToken(t, userID)}
		req.AddCookie(cookie)
	}

	csrf := avoidCSRF{}
	csrf.SetCSRFCookie(req)
	csrf.SetCSRFHeader(req)

	rec := httptest.NewRecorder()
	return req, rec
}

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
				SetUpdatedAt(d.CreatedAt).
				SetUser(user).
				SaveX(context.Background())
		}

		req, rec := createAuthenticatedRequest(t, http.MethodGet, "/todo", "", user.ID)
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res dto.ListTodoResponseDto
		json.Unmarshal(rec.Body.Bytes(), &res)
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
}

func TestTodoHandler_ListTodo_IncludeDone_Integration(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

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
		json.Unmarshal(rec.Body.Bytes(), &res)
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
		json.Unmarshal(rec.Body.Bytes(), &res)
		assert.Len(t, res.Data, 3)
	})

	t.Run("include_done=false returns only active todos", func(t *testing.T) {
		req, rec := createAuthenticatedRequest(t, http.MethodGet, "/todo?include_done=false", "", user.ID)
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res dto.ListTodoResponseDto
		json.Unmarshal(rec.Body.Bytes(), &res)
		assert.Len(t, res.Data, 2)
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

		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())
		body := `{"title": "New Todo", "description": "New Description"}`

		req, rec := createAuthenticatedRequest(t, http.MethodPost, "/todo", body, user.ID)
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
		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

		body := `{"title": "", "description": "New Description"}`

		req, rec := createAuthenticatedRequest(t, http.MethodPost, "/todo", body, user.ID)
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expected := `{"error":{"title":"タイトルは必須です"}}`
		assert.JSONEq(t, expected, rec.Body.String())
	})

	t.Run("Todo 新規作成 未認証", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

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

		req, rec := createAuthenticatedRequest(t, http.MethodPatch, fmt.Sprintf("/todo/%d", todo.ID), body, user.ID)
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
		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

		body := `{"title": "Updated Title"}`

		req, rec := createAuthenticatedRequest(t, http.MethodPatch, "/todo/999", body, user.ID)
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Todo 更新 他人のTodo", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

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
	t.Run("Todo 更新 完了済み", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

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

		req, rec := createAuthenticatedRequest(t, http.MethodDelete, fmt.Sprintf("/todo/%d", todo.ID), "", user.ID)
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
		user := client.User.Create().SetName("test").SetEmail("test").SetPassword("test").SaveX(context.Background())

		req, rec := createAuthenticatedRequest(t, http.MethodDelete, "/todo/999", "", user.ID)
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Todo 削除 他人のTodo", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer client.Close()

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
}

func TestTodoHandler_UpdateDoneStatus(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	t.Run("成功時", func(t *testing.T) {
		e := echo.New()
		reqBody := map[string]interface{}{"is_done": true}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/todo/1/done", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/todo/:id/done")
		c.SetPathValues(echo.PathValues{
			{Name: "id", Value: "1"},
		})

		repo := &spyTodoRepo{
			updateDoneStatus: func(id int, isDone bool) (*ent.Todo, error) {
				now := time.Now()
				return &ent.Todo{
					ID:        id,
					Title:     "Todo",
					DoneAt:    &now,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}

		factory := func(logger *slog.Logger) (*services.TodoService, error) {
			return services.NewTodoService(logger, repo), nil
		}

		h := handlers.NewTodoHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), client, factory)

		err := h.UpdateDoneStatus(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "\"done_at\"")
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		e := echo.New()
		reqBody := map[string]interface{}{} // is_done missing
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/todo/1/done", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/todo/:id/done")
		c.SetPathValues(echo.PathValues{
			{Name: "id", Value: "1"},
		})

		h := handlers.NewTodoHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), client, nil)

		err := h.UpdateDoneStatus(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "error")
	})

	t.Run("Todoが見つからない場合", func(t *testing.T) {
		e := echo.New()
		reqBody := map[string]interface{}{"is_done": true}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/todo/999/done", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/todo/:id/done")
		c.SetPathValues(echo.PathValues{
			{Name: "id", Value: "999"},
		})

		repo := &spyTodoRepo{
			updateDoneStatus: func(id int, isDone bool) (*ent.Todo, error) {
				return nil, &ent.NotFoundError{}
			},
		}

		factory := func(logger *slog.Logger) (*services.TodoService, error) {
			return services.NewTodoService(logger, repo), nil
		}

		h := handlers.NewTodoHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), client, factory)

		err := h.UpdateDoneStatus(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
