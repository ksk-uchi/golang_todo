//go:build wireinject
// +build wireinject

package main

import (
	"todo-app/handlers"
	"todo-app/providers"
	"todo-app/routes"
	"todo-app/services"

	"github.com/google/wire"
	"github.com/labstack/echo/v5"
)

// todo
var todoSet = wire.NewSet(
	handlers.NewTodoHandler,
	routes.NewTodoRouter,
	services.ProvideTodoServiceFactory,
)

// app
var appSet = wire.NewSet(
	providers.NewEntClient,
	routes.NewRouter,
	echo.New,
	NewApp,
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
		todoSet,
		appSet,
	)
	return &App{}, nil, nil
}
