package services

import (
	"context"
	"testing"

	"todo-app/dto"
	"todo-app/ent"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ent.User), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		req           *dto.LoginRequest
		mockUser      *ent.User
		mockError     error
		expectedToken bool
		expectedError string
	}{
		{
			name: "Success",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: password,
			},
			mockUser: &ent.User{
				ID:       1,
				Email:    "test@example.com",
				Password: string(hashedPassword),
			},
			mockError:     nil,
			expectedToken: true,
			expectedError: "",
		},
		{
			name: "UserNotFound",
			req: &dto.LoginRequest{
				Email:    "notfound@example.com",
				Password: password,
			},
			mockUser:      nil,
			mockError:     &ent.NotFoundError{},
			expectedToken: false,
			expectedError: "invalid email or password",
		},
		{
			name: "InvalidPassword",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockUser: &ent.User{
				ID:       1,
				Email:    "test@example.com",
				Password: string(hashedPassword),
			},
			mockError:     nil,
			expectedToken: false,
			expectedError: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			authService := NewAuthService(mockRepo)

			mockRepo.On("FindByEmail", mock.Anything, tt.req.Email).Return(tt.mockUser, tt.mockError)

			token, err := authService.Login(context.Background(), tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
