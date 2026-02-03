package handlers

import (
	"net/http"
	"todo-app/ent"
	"todo-app/ent/todo"
	"todo-app/middlewares" // トランザクション取得用
	"todo-app/services"

	"github.com/labstack/echo/v5"
)

type TodoHandler struct {
	service services.ITodoService
}

func NewTodoHandler(s services.ITodoService) *TodoHandler {
	return &TodoHandler{service: s}
}

func (h *TodoHandler) ListTodo(c *echo.Context) error {
	tx := middlewares.GetTx(c)

	ctx := c.Request().Context()
	todos, err := tx.Todo.Query().
		Order(ent.Desc(todo.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch todos",
		})
	}

	res := h.service.EntitiesToDTOs(todos)

	return c.JSON(http.StatusOK, res)
}
