package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"todo-app/dto"
	"todo-app/services"
	"todo-app/utils"
	"todo-app/validators"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	logger  *slog.Logger
	service services.IAuthService
}

func NewAuthHandler(logger *slog.Logger, service services.IAuthService) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: service,
	}
}

func (h *AuthHandler) Login(c *echo.Context) error {
	utils.LogRequest(h.logger, c)

	var req validators.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.HandleError(h.logger, c, errors.New("invalid request"), http.StatusBadRequest)
	}

	if errs := req.Validate(); errs != nil {
		h.logger.Error("validation error", slog.Any("errors", errs))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": errs})
	}

	input := &dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := h.service.Login(c.Request().Context(), input)
	if err != nil {
		return utils.HandleError(h.logger, c, errors.New("invalid email or password"), http.StatusUnauthorized)
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "login success"})
}
