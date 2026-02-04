package handlers

import (
	"log/slog"
	"net/http"
	"strings"
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

type createTodoRequest struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=200"`
}

func (h *TodoHandler) CreateTodo(c *echo.Context) error {
	errorHandling := func(c *echo.Context, err error, code int) error {
		h.logger.Error(err.Error())
		return c.JSON(code, map[string]string{
			"error": err.Error(),
		})
	}

	var req createTodoRequest
	if err := c.Bind(&req); err != nil {
		return errorHandling(c, err, http.StatusBadRequest)
	}

	if err := validate.Struct(&req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return errorHandling(c, err, http.StatusBadRequest)
		}

		errorMessages := make(map[string]string)
		for _, fe := range validationErrors {
			field := strings.ToLower(fe.Field())
			switch field {
			case "title":
				switch fe.Tag() {
				case "required":
					errorMessages[field] = "タイトルは必須です"
				case "max":
					errorMessages[field] = "タイトルは" + fe.Param() + "文字以内で入力してください"
				}
			case "description":
				switch fe.Tag() {
				case "max":
					errorMessages[field] = "説明は" + fe.Param() + "文字以内で入力してください"
				}
			}
		}

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
