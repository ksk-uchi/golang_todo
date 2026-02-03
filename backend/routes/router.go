package routes

import (
	"github.com/labstack/echo/v5"
)

func NewRouter(txM echo.MiddlewareFunc, todoR *TodoRouter) *Router {
	return &Router{
		txMiddleware: txM,
		todo:         todoR,
	}
}

type Router struct {
	txMiddleware echo.MiddlewareFunc
	todo         *TodoRouter
}

func (r *Router) Setup(e *echo.Echo) {
	e.Use(r.txMiddleware)

	r.todo.SetupTodoRoute(e.Group("/todo"))
}
