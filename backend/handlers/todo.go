package handlers

import (
	"net/http" // トランザクション取得用
	"todo-app/dto"
	"todo-app/ent"
	"todo-app/services"

	"github.com/labstack/echo/v5"
)

type TodoHandler struct {
	client         *ent.Client
	serviceFactory services.TodoServiceFactory
}

func NewTodoHandler(client *ent.Client, factory services.TodoServiceFactory) *TodoHandler {
	return &TodoHandler{
		client:         client,
		serviceFactory: factory,
	}
}

func (h *TodoHandler) ListTodo(c *echo.Context) error {
	ctx := c.Request().Context()

	service := h.serviceFactory(h.client)
	todos, err := service.GetTodoSlice(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch todos",
		})
	}

	res := dto.EntitiesToTodoDtoSlice(todos)

	return c.JSON(http.StatusOK, res)
}
