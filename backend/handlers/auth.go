package handlers

import (
	"net/http"

	"todo-app/dto"
	"todo-app/services"
	"todo-app/validators"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	service services.IAuthService
}

func NewAuthHandler(service services.IAuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var req validators.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if errs := req.Validate(); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": errs})
	}

	input := &dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := h.service.Login(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
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
