//go:build wireinject
// +build wireinject

package di

import (
	"todo-app/ent"
	"todo-app/handlers"
	"todo-app/middleware"
	"todo-app/providers"
	"todo-app/repositories"
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

// auth
var authSet = wire.NewSet(
	repositories.NewUserRepository,
	wire.Bind(new(repositories.IUserRepository), new(*repositories.UserRepository)),
	services.NewAuthService,
	wire.Bind(new(services.IAuthService), new(*services.AuthService)),
	handlers.NewAuthHandler,
	routes.NewAuthRouter,
	middleware.NewAuthMiddleware,
)

// app
var appSet = wire.NewSet(
	providers.NewEntClient,
	routes.NewRouter,
	NewLogger,
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
		authSet,
		appSet,
	)
	return &App{}, nil, nil
}

func InitializeTestApp(e *echo.Echo, client *ent.Client) (*App, error) {
	wire.Build(
		todoSet,
		authSet,
		routes.NewRouter,
		NewLogger,
		NewApp,
	)
	return &App{}, nil
}
