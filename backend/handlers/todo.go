package handlers

import (
	"log/slog"
	"net/http"
	"todo-app/dto"
	"todo-app/ent"
	"todo-app/services"
	"todo-app/validators"

	"github.com/labstack/echo/v5"
)

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

	var req validators.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return errorHandling(c, err, http.StatusBadRequest)
	}

	if errorMessages := req.Validate(); errorMessages != nil {
		h.logger.Error("validation error", slog.Any("errors", errorMessages))
		return c.JSON(http.StatusBadRequest, map[string]map[string]string{
			"error": errorMessages,
		})
	}

	ctx := c.Request().Context()
	service, err := h.serviceFactory(ctx, h.logger, h.client)
	if err != nil {
		return errorHandling(c, err, http.StatusInternalServerError)
	}

	todo, err := service.CreateTodo(req.Title, req.Description)
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

func (h *TodoHandler) UpdateTodo(c *echo.Context) error {
	errorHandling := func(c *echo.Context, err error, code int) error {
		h.logger.Error(err.Error())
		return c.JSON(code, map[string]string{
			"error": err.Error(),
		})
	}

	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idParam"})
	}

	var req validators.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return errorHandling(c, err, http.StatusBadRequest)
	}

	if errorMessages := req.Validate(); errorMessages != nil {
		h.logger.Error("validation error", slog.Any("errors", errorMessages))
		return c.JSON(http.StatusBadRequest, map[string]map[string]string{
			"error": errorMessages,
		})
	}

	ctx := c.Request().Context()
	service, err := h.serviceFactory(ctx, h.logger, h.client)
	if err != nil {
		return errorHandling(c, err, http.StatusInternalServerError)
	}

	todo, err := service.UpdateTodo(id, req.Title, req.Description)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "todo not found"})
		}
		return errorHandling(c, err, http.StatusInternalServerError)
	}

	res := dto.EntityToTodoDto(todo)

	return c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) DeleteTodo(c *echo.Context) error {
	errorHandling := func(c *echo.Context, err error, code int) error {
		h.logger.Error(err.Error())
		return c.JSON(code, map[string]string{
			"error": err.Error(),
		})
	}

	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idParam"})
	}

	ctx := c.Request().Context()
	service, err := h.serviceFactory(ctx, h.logger, h.client)
	if err != nil {
		return errorHandling(c, err, http.StatusInternalServerError)
	}

	err = service.DeleteTodo(id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.NoContent(http.StatusNoContent)
		}
		return errorHandling(c, err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
