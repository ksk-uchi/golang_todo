package routes

import (
	"github.com/labstack/echo/v5"
)

func (r *Router) SetupTodo(eg *echo.Group) {
	eg.GET("", r.TodoHandler.ListTodo)
}
