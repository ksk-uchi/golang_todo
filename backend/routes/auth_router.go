package routes

import (
	"todo-app/handlers"

	"github.com/labstack/echo/v5"
)

type AuthRouter struct {
	handler *handlers.AuthHandler
}

func NewAuthRouter(handler *handlers.AuthHandler) *AuthRouter {
	return &AuthRouter{handler: handler}
}

func (r *AuthRouter) SetupAuthRoute(g *echo.Group) {
	g.POST("/login", r.handler.Login)
}
