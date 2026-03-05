package services

import (
	"context"
	"testing"
	"time"
	"todo-app/ent"
	"todo-app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genai"
)

type MockGenAIClient struct {
	mock.Mock
}

func (m *MockGenAIClient) GenerateContent(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	args := m.Called(ctx, model, contents, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*genai.GenerateContentResponse), args.Error(1)
}

func TestDecideFilterTodosFunction(t *testing.T) {
	oldCreateAIClientFunc := createAIClientFunc
	defer func() { createAIClientFunc = oldCreateAIClientFunc }()

	mockClient := new(MockGenAIClient)
	createAIClientFunc = func(ctx context.Context) (IGenAIClient, error) {
		return mockClient, nil
	}

	repo := new(testutils.MockTodoRepository)
	s := NewAIService(repo)
	ctx := context.Background()

	t.Run("success - function call", func(t *testing.T) {
		expectedResponse := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []*genai.Part{
							{
								FunctionCall: &genai.FunctionCall{
									Name: "ListTodosByDoneAt",
									Args: map[string]interface{}{
										"done_from": "2023-01-01T00:00:00Z",
										"done_to":   "2023-01-31T23:59:59Z",
									},
								},
							},
						},
					},
				},
			},
		}
		mockClient.On("GenerateContent", ctx, "gemini-3-flash-preview", mock.Anything, mock.Anything).Return(expectedResponse, nil).Once()

		res, err := s.DecideFilterTodosFunction(ctx, "last month")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "ListTodosByDoneAt", res.FunctionName)
		assert.Equal(t, "2023-01-01T00:00:00Z", res.Args["done_from"])
	})

	t.Run("success - no function call", func(t *testing.T) {
		expectedResponse := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []*genai.Part{
							{
								Text: "対応できる Tool がありません",
							},
						},
					},
				},
			},
		}
		mockClient.On("GenerateContent", ctx, "gemini-3-flash-preview", mock.Anything, mock.Anything).Return(expectedResponse, nil).Once()

		res, err := s.DecideFilterTodosFunction(ctx, "hello")
		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

func TestFilterTodos(t *testing.T) {
	repo := new(testutils.MockTodoRepository)
	s := NewAIService(repo)
	ctx := context.Background()

	t.Run("ListTodosByDoneAt - success", func(t *testing.T) {
		args := map[string]interface{}{
			"done_from": "2023-01-01T00:00:00Z",
			"done_to":   "2023-01-31T23:59:59Z",
		}

		doneFrom, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
		doneTo, _ := time.Parse(time.RFC3339, "2023-01-31T23:59:59Z")
		repo.On("FetchTodosByDoneAt", mock.Anything, &doneFrom, &doneTo).Return([]*ent.Todo{{ID: 1, Title: "Test"}}, nil)

		res, err := s.FilterTodos(ctx, "ListTodosByDoneAt", args)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 1, res[0].ID)
	})

	t.Run("Unknown function", func(t *testing.T) {
		res, err := s.FilterTodos(ctx, "Unknown", nil)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
