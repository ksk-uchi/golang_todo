package routes

import (
	"github.com/labstack/echo/v5"
)

func NewRouter(todoR *TodoRouter) *Router {
	return &Router{
		todo: todoR,
	}
}

type Router struct {
	todo *TodoRouter
}

func (r *Router) Setup(e *echo.Echo) {
	r.todo.SetupTodoRoute(e.Group("/todo"))
}
