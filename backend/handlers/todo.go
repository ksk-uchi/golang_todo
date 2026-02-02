package handlers

import (
	"net/http"
	"todo-app/middlewares"
	"todo-app/services"

	"github.com/labstack/echo/v5"
)

func (s *Handler) ListTodo(c *echo.Context) error {
	tx := middlewares.GetTx(c)
	todos, err := tx.Todo.Query().All(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	service := services.NewTodoService()
	res := service.EntitiesToDTOs(todos)

	return c.JSON(http.StatusOK, res)
}
