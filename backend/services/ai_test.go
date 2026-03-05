package services

import (
	"context"
	"testing"
	"time"
	"todo-app/ent"

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

type spyTodoRepo struct {
	fetchTodosByDoneAt func(doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error)
}

func (s *spyTodoRepo) FetchTodos(ctx context.Context, limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
	return nil, nil
}
func (s *spyTodoRepo) GetTodoCount(ctx context.Context, includeDone bool) (int, error) {
	return 0, nil
}
func (s *spyTodoRepo) FindTodo(ctx context.Context, id int) (*ent.Todo, error) {
	return nil, nil
}
func (s *spyTodoRepo) GetTodoForUpdate(ctx context.Context, id int) (*ent.Todo, error) {
	return nil, nil
}
func (s *spyTodoRepo) CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error) {
	return nil, nil
}
func (s *spyTodoRepo) UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error) {
	return nil, nil
}
func (s *spyTodoRepo) UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error) {
	return nil, nil
}
func (s *spyTodoRepo) DeleteTodo(ctx context.Context, id int) error {
	return nil
}
func (s *spyTodoRepo) FetchTodosByDoneAt(ctx context.Context, doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error) {
	if s.fetchTodosByDoneAt != nil {
		return s.fetchTodosByDoneAt(doneFrom, doneTo)
	}
	return nil, nil
}

func TestDecideFilterTodosFunction(t *testing.T) {
	oldCreateAIClientFunc := createAIClientFunc
	defer func() { createAIClientFunc = oldCreateAIClientFunc }()

	mockClient := new(MockGenAIClient)
	createAIClientFunc = func(ctx context.Context) (IGenAIClient, error) {
		return mockClient, nil
	}

	repo := &spyTodoRepo{}
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
	repo := &spyTodoRepo{}
	s := NewAIService(repo)
	ctx := context.Background()

	t.Run("ListTodosByDoneAt - success", func(t *testing.T) {
		args := map[string]interface{}{
			"done_from": "2023-01-01T00:00:00Z",
			"done_to":   "2023-01-31T23:59:59Z",
		}

		repo.fetchTodosByDoneAt = func(doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error) {
			assert.Equal(t, "2023-01-01T00:00:00Z", doneFrom.Format(time.RFC3339))
			assert.Equal(t, "2023-01-31T23:59:59Z", doneTo.Format(time.RFC3339))
			return []*ent.Todo{{ID: 1, Title: "Test"}}, nil
		}

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
