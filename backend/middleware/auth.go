package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"todo-app/ent"
	"todo-app/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type AuthMiddleware struct {
	userRepo repositories.IUserRepository
}

func NewAuthMiddleware(userRepo repositories.IUserRepository) *AuthMiddleware {
	return &AuthMiddleware{userRepo: userRepo}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Skip authentication for login and options requests
		if c.Path() == "/auth/login" || c.Request().Method == http.MethodOptions {
			return next(c)
		}

		cookie, err := c.Cookie("token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
		}

		tokenString := cookie.Value
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secret"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
		}
		userID := int(userIDFloat)

		user, err := m.userRepo.FindById(c.Request().Context(), userID)
		if err != nil {
			if ent.IsNotFound(err) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}

		// Update token expiration
		newClaims := jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		newTokenString, err := newToken.SignedString([]byte(secret))
		if err == nil {
			newCookie := new(http.Cookie)
			newCookie.Name = "token"
			newCookie.Value = newTokenString
			newCookie.Path = "/"
			newCookie.HttpOnly = true
			newCookie.SameSite = http.SameSiteStrictMode
			newCookie.Expires = time.Now().Add(time.Hour)
			c.SetCookie(newCookie)
		}

		ctx := context.WithValue(c.Request().Context(), "user", user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
