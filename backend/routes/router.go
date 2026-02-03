package routes

import (
	"todo-app/ent"
	"todo-app/handlers"

	"github.com/labstack/echo/v5"
)

func NewRouter(todoH *handlers.TodoHandler, client *ent.Client, txM echo.MiddlewareFunc) *Router {
	return &Router{
		client:       client,
		TodoHandler:  todoH,
		txMiddleware: txM,
	}
}

type Router struct {
	client       *ent.Client
	TodoHandler  *handlers.TodoHandler
	txMiddleware echo.MiddlewareFunc
}

func (r *Router) Setup(e *echo.Echo) {
	e.Use(r.txMiddleware)

	r.SetupTodo(e.Group("/todo"))
}
