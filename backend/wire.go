//go:build wireinject
// +build wireinject

package main

import (
	"todo-app/handlers"
	"todo-app/middlewares"
	"todo-app/providers"
	"todo-app/routes"
	"todo-app/services"

	"github.com/google/wire"
	"github.com/labstack/echo/v5"
)

type App struct {
	Engine *echo.Echo
	Router *routes.Router
}

func NewApp(e *echo.Echo, r *routes.Router) *App {
	return &App{Engine: e, Router: r}
}

func InitializeApp() (*App, func(), error) {
	wire.Build(
		middlewares.NewTransactionMiddleware,
		providers.NewEntClient,
		services.NewTodoService,
		handlers.NewTodoHandler,
		routes.NewRouter,
		echo.New,
		NewApp,
	)
	return &App{}, nil, nil
}
