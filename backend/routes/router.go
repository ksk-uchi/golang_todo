package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-CSRF-Token",
		},
		AllowCredentials: true,
	}))
	r.todo.SetupTodoRoute(e.Group("/todo"))
}
