package main

import (
	"log/slog"

	"todo-app/ent"
	"todo-app/handlers"
	"todo-app/middlewares"
	"todo-app/routes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	client, err := ent.Open("mysql", "user:password@tcp(localhost:3306)/todo_db?parseTime=True")
	if err != nil {
		slog.Error("failed opening connection to mysql.", "error", err)
	}
	defer client.Close()

	// Middleware
	e.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	e.Use(middleware.Recover())       // recover panics as errors for proper error handling
	e.Use(middlewares.TransactionMiddleware(client))

	h := handlers.NewHandler(client)

	routes.InitRoutes(e, h)

	// Start server
	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
