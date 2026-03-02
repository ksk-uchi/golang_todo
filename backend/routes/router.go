package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"todo-app/middleware"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

func NewRouter(todoR *TodoRouter, authR *AuthRouter, authM *middleware.AuthMiddleware) *Router {
	return &Router{
		todo:  todoR,
		auth:  authR,
		authM: authM,
	}
}

type Router struct {
	todo  *TodoRouter
	auth  *AuthRouter
	authM *middleware.AuthMiddleware
}

func (r *Router) Setup(e *echo.Echo) {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		origins = os.Getenv("FRONTEND_ORIGIN")
	}
	if origins == "" {
		origins = "http://localhost:3000"
	}
	allowOrigins := strings.Split(origins, ",")

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-CSRF-Token",
		},
		AllowCredentials: true,
	}))

	e.Use(r.authM.Authenticate)
	skipper := func(c *echo.Context) bool {
		return c.Request().URL.Path == "/auth/login"
	}
	e.Use(echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		Skipper:        skipper,
		TrustedOrigins: []string{os.Getenv("FRONTEND_ORIGIN")},
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookiePath:     "/",
		CookieSameSite: http.SameSiteStrictMode,
	}))

	// CSP Nonce の動的生成 middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			nonce := make([]byte, 16)
			if _, err := rand.Read(nonce); err != nil {
				return err
			}
			nonceStr := base64.StdEncoding.EncodeToString(nonce)
			c.Set("nonce", nonceStr)
			return next(c)
		}
	})

	// XSS 対策 (Secure 内部の CSP はここでは使わず、個別に Header を設定)
	e.Use(echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
		// Static な CSP は空にし、middleware で動的に設定する
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			nonceVal := c.Get("nonce")
			if nonceVal == nil {
				return next(c)
			}
			nonce := nonceVal.(string)
			csp := fmt.Sprintf("default-src 'self'; script-src 'self' 'nonce-%s'; object-src 'none'; base-uri 'self';", nonce)
			c.Response().Header().Set(echo.HeaderContentSecurityPolicy, csp)
			return next(c)
		}
	})

	r.auth.SetupAuthRoute(e.Group("/auth"))
	r.todo.SetupTodoRoute(e.Group("/todo"))
}
