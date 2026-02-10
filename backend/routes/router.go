package routes

import (
	"net/http"
	"todo-app/middleware"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

func NewRouter(todoR *TodoRouter, authR *AuthRouter, authM *middleware.AuthMiddleware) *Router {
	return &Router{
		todo:  todoR,
		auth:  authR,
		authM: authM,
	}
}

type Router struct {
	todo  *TodoRouter
	auth  *AuthRouter
	authM *middleware.AuthMiddleware
}

func (r *Router) Setup(e *echo.Echo) {
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
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

	r.auth.SetupAuthRoute(e.Group("/auth"))

	e.Use(r.authM.Authenticate)
	skipper := func(c *echo.Context) bool {
		return c.Request().URL.Path == "/auth/login"
	}
	e.Use(echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		Skipper:        skipper,
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookiePath:     "/",
		CookieSameSite: http.SameSiteStrictMode,
	}))
	r.todo.SetupTodoRoute(e.Group("/todo"))
}
