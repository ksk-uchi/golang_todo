package routes

import (
	"todo-app/handlers"

	"github.com/labstack/echo/v5"
)

func NewTodoRouter(todoH *handlers.TodoHandler) *TodoRouter {
	return &TodoRouter{
		TodoHandler: todoH,
	}
}

type TodoRouter struct {
	TodoHandler *handlers.TodoHandler
}

func (r *TodoRouter) SetupTodoRoute(eg *echo.Group) {
	eg.GET("", r.TodoHandler.ListTodo)
}
