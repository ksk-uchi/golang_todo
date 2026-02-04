package handlers

import (
	"log/slog"
	"net/http"
	"todo-app/dto"
	"todo-app/ent"
	"todo-app/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var validate = validator.New()

type TodoHandler struct {
	logger         *slog.Logger
	client         *ent.Client
	serviceFactory services.TodoServiceFactory
}

func NewTodoHandler(logger *slog.Logger, client *ent.Client, factory services.TodoServiceFactory) *TodoHandler {
	return &TodoHandler{
		logger:         logger,
		client:         client,
		serviceFactory: factory,
	}
}

func (h *TodoHandler) CreateTodo(c *echo.Context) error {
	errorHandling := func(c *echo.Context, err error, code int) error {
		h.logger.Error(err.Error())
		return c.JSON(code, map[string]string{
			"error": err.Error(),
		})
	}

	var createTodoDto dto.CreateTodoDto
	if err := c.Bind(&createTodoDto); err != nil {
		return errorHandling(c, err, http.StatusBadRequest)
	}

	if err := validate.Struct(&createTodoDto); err != nil {
		return errorHandling(c, err, http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	service, err := h.serviceFactory(ctx, h.logger, h.client)
	if err != nil {
		return errorHandling(c, err, http.StatusInternalServerError)
	}

	todo, err := service.CreateTodo(&createTodoDto)
	if err != nil {
		return errorHandling(c, err, http.StatusInternalServerError)
	}

	res := dto.EntityToTodoDto(todo)

	return c.JSON(http.StatusCreated, res)
}

func (h *TodoHandler) ListTodo(c *echo.Context) error {
	errorHandling := func(c *echo.Context, err error) error {
		h.logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch todos",
		})
	}

	ctx := c.Request().Context()
	service, err := h.serviceFactory(ctx, h.logger, h.client)
	if err != nil {
		return errorHandling(c, err)
	}
	todos, err := service.GetTodoSlice()
	if err != nil {
		return errorHandling(c, err)
	}

	res := dto.EntitiesToTodoDtoSlice(todos)

	return c.JSON(http.StatusOK, res)
}
