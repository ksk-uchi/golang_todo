package routes

import (
	"todo-app/handlers"

	"github.com/labstack/echo/v5"
)

// InitRoutes に echo インスタンスと依存関係(Server)を渡します
func InitRoutes(e *echo.Echo, h *handlers.Handler) {
	// Todo関連
	e.GET("/todo", h.ListTodo)
}
