package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"todo-app/dto"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		reqBody        string
		mockToken      string
		mockError      error
		expectedStatus int
		expectedCookie bool
	}{
		{
			name:           "Success",
			reqBody:        `{"email":"test@example.com","password":"password"}`,
			mockToken:      "mock_token",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCookie: true,
		},
		{
			name:           "InvalidRequest",
			reqBody:        `invalid_json`,
			mockToken:      "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCookie: false,
		},
		{
			name:           "ServiceError",
			reqBody:        `{"email":"test@example.com","password":"password"}`,
			mockToken:      "",
			mockError:      errors.New("invalid email or password"),
			expectedStatus: http.StatusUnauthorized,
			expectedCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			handler := NewAuthHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.name != "InvalidRequest" {
				var reqDTO dto.LoginRequest
				_ = json.Unmarshal([]byte(tt.reqBody), &reqDTO)
				mockService.On("Login", mock.Anything, mock.MatchedBy(func(r *dto.LoginRequest) bool {
					return r.Email == reqDTO.Email && r.Password == reqDTO.Password
				})).Return(tt.mockToken, tt.mockError)
			}

			err := handler.Login(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedCookie {
				cookies := rec.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "token" {
						assert.Equal(t, tt.mockToken, cookie.Value)
						assert.True(t, cookie.HttpOnly)
						assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
						assert.WithinDuration(t, time.Now().Add(time.Hour), cookie.Expires, time.Minute)
						found = true
						break
					}
				}
				assert.True(t, found, "cookie 'token' should be set")
			}
		})
	}
}
