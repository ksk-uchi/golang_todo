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

// provider
var providerSet = wire.NewSet(
	providers.NewEntClient,
)

// middleware
var middlewareSet = wire.NewSet(
	middlewares.NewTransactionMiddleware,
)

// handler
var handlerSet = wire.NewSet(
	handlers.NewTodoHandler,
)

// service
var serviceSet = wire.NewSet(
	services.NewTodoService,
)

// route
var routeSet = wire.NewSet(
	routes.NewRouter,
	routes.NewTodoRouter,
)

// app
var appSet = wire.NewSet(
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
		providerSet,
		middlewareSet,
		handlerSet,
		serviceSet,
		routeSet,
		appSet,
	)
	return &App{}, nil, nil
}
