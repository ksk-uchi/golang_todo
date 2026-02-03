package handlers

import (
	"net/http" // トランザクション取得用
	"todo-app/dto"
	"todo-app/services"

	"github.com/labstack/echo/v5"
)

func NewTodoHandler(s *services.TodoService) *TodoHandler {
	return &TodoHandler{service: s}
}

type TodoHandler struct {
	service *services.TodoService
}

func (h *TodoHandler) ListTodo(c *echo.Context) error {
	todos, err := h.service.GetTodoSlice(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch todos",
		})
	}
	res := dto.EntitiesToTodoDtoSlice(todos)

	return c.JSON(http.StatusOK, res)
}
