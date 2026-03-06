package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"todo-app/app_errors"
	"todo-app/dto"
	"todo-app/ent"
	"todo-app/services"
	"todo-app/utils"
	"todo-app/validators"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type TodoHandler struct {
	logger               *slog.Logger
	service              *services.TodoService
	filterHistoryService services.ITodoFilterHistoryService
	aiService            *services.AIService
}

func NewTodoHandler(logger *slog.Logger, service *services.TodoService, filterHistoryService services.ITodoFilterHistoryService, aiService *services.AIService) *TodoHandler {
	return &TodoHandler{
		logger:               logger,
		service:              service,
		filterHistoryService: filterHistoryService,
		aiService:            aiService,
	}
}

func (h *TodoHandler) ListTodoFilterHistories(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	ctx := c.Request().Context()
	histories, err := h.filterHistoryService.FetchLatestFilters(ctx)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	res := dto.ListTodoFilterHistoriesResponseDto{
		Queries: dto.EntitiesToTodoFilterHistoryQueryDtos(histories),
	}

	return c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) FilterTodosByQuery(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	query := c.QueryParam("query")
	ctx := c.Request().Context()

	aiDto, err := h.aiService.DecideFilterTodosFunction(ctx, query)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	var todos []*ent.Todo
	var functionName *string
	var args map[string]interface{}

	if aiDto != nil {
		functionName = &aiDto.FunctionName
		args = aiDto.Args
		todos, err = h.aiService.FilterTodos(ctx, aiDto.FunctionName, aiDto.Args)
		if err != nil {
			return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
		}
	} else {
		todos = []*ent.Todo{}
	}

	todoIds := make([]int, len(todos))
	for i, t := range todos {
		todoIds[i] = t.ID
	}

	_, err = h.filterHistoryService.SaveFilterHistory(ctx, query, functionName, args, todoIds)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	res := make([]dto.TodoDto, len(todos))
	for i, t := range todos {
		res[i] = dto.EntityToTodoDto(t)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) FilterTodosByQueryID(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	queryIDStr := c.QueryParam("query_id")
	queryID, err := uuid.Parse(queryIDStr)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	history, err := h.filterHistoryService.GetFilterHistoryByQueryID(ctx, queryID)
	if err != nil {
		if ent.IsNotFound(err) {
			return utils.HandleError(h.logger, c, err, http.StatusNotFound)
		}
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	todos, err := h.service.FetchTodosByIds(ctx, history.ResultTodoIds)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	res := make([]dto.TodoDto, len(todos))
	for i, t := range todos {
		res[i] = dto.EntityToTodoDto(t)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) CreateTodo(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	var req validators.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusBadRequest)
	}

	if errorMessages := req.Validate(); errorMessages != nil {
		h.logger.Error("validation error", slog.Any("errors", errorMessages))
		return c.JSON(http.StatusBadRequest, map[string]map[string]string{
			"error": errorMessages,
		})
	}

	ctx := c.Request().Context()
	todo, err := h.service.CreateTodo(ctx, req.Title, req.Description)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	res := dto.EntityToTodoDto(todo)

	return c.JSON(http.StatusCreated, res)
}

func (h *TodoHandler) ListTodo(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	pageInt, err := echo.QueryParamOr(c, "page", 1)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := echo.QueryParamOr(c, "limit", 20)
	if err != nil || limitInt < 1 {
		limitInt = 20
	}
	if limitInt > 100 {
		limitInt = 100
	}

	includeDone, err := echo.QueryParamOr(c, "include_done", false)
	if err != nil {
		includeDone = false
	}

	ctx := c.Request().Context()
	todos, err := h.service.GetTodoSlice(ctx, pageInt, limitInt, includeDone)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	pagination, err := h.service.CalculatePagination(ctx, pageInt, limitInt, includeDone)
	if err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	res := dto.ListTodoResponseDto{
		Data:       todos,
		Pagination: pagination,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) UpdateTodo(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return utils.HandleError(h.logger, c, errors.New("invalid idParam"), http.StatusBadRequest)
	}

	var req validators.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusBadRequest)
	}

	if errorMessages := req.Validate(); errorMessages != nil {
		h.logger.Error("validation error", slog.Any("errors", errorMessages))
		return c.JSON(http.StatusBadRequest, map[string]map[string]string{
			"error": errorMessages,
		})
	}

	ctx := c.Request().Context()
	todo, err := h.service.UpdateTodo(ctx, id, req.Title, req.Description)
	if err != nil {
		if ent.IsNotFound(err) {
			return utils.HandleError(h.logger, c, errors.New("todo not found"), http.StatusNotFound)
		}
		if errors.Is(err, app_errors.ErrTodoAlreadyDone) {
			return utils.HandleError(h.logger, c, err, http.StatusBadRequest)
		}
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	res := dto.EntityToTodoDto(todo)

	return c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) DeleteTodo(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return utils.HandleError(h.logger, c, errors.New("invalid idParam"), http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	err = h.service.DeleteTodo(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.NoContent(http.StatusNoContent)
		}
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TodoHandler) UpdateDoneStatus(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return utils.HandleError(h.logger, c, errors.New("invalid idParam"), http.StatusBadRequest)
	}

	var req validators.UpdateDoneStatusRequest
	if err := c.Bind(&req); err != nil {
		return utils.HandleError(h.logger, c, err, http.StatusBadRequest)
	}

	if errorMessages := req.Validate(); errorMessages != nil {
		h.logger.Error("validation error", slog.Any("errors", errorMessages))
		return c.JSON(http.StatusBadRequest, map[string]map[string]string{
			"error": errorMessages,
		})
	}

	ctx := c.Request().Context()
	todo, err := h.service.UpdateDoneStatus(ctx, id, *req.IsDone)
	if err != nil {
		if ent.IsNotFound(err) {
			return utils.HandleError(h.logger, c, errors.New("todo not found"), http.StatusNotFound)
		}
		return utils.HandleError(h.logger, c, err, http.StatusInternalServerError)
	}

	res := dto.EntityToTodoDto(todo)
	return c.JSON(http.StatusOK, res)
}
