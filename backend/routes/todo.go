package routes

import (
	"todo-app/handlers"

	"github.com/labstack/echo/v5"
)

func NewTodoRouter(todoH *handlers.TodoHandler, txM echo.MiddlewareFunc) *TodoRouter {
	return &TodoRouter{
		TodoHandler:  todoH,
		txMiddleware: txM,
	}
}

type TodoRouter struct {
	TodoHandler  *handlers.TodoHandler
	txMiddleware echo.MiddlewareFunc
}

func (r *TodoRouter) SetupTodoRoute(eg *echo.Group) {
	eg.GET("", r.TodoHandler.ListTodo)
}
