package services_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"todo-app/ent"
	"todo-app/services"
	"todo-app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTodoFilterHistoryService_FetchLatestFilters(t *testing.T) {
	t.Run("リポジトリを呼び出し、最新5件の履歴を返すこと", func(t *testing.T) {
		repo := new(testutils.MockTodoFilterHistoryRepository)
		expectedHistories := []*ent.TodoFilterHistory{
			{Query: "query 1"},
			{Query: "query 2"},
			{Query: "query 3"},
			{Query: "query 4"},
			{Query: "query 5"},
		}
		repo.On("FetchLatestFilters", mock.Anything, 5).Return(expectedHistories, nil)

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := services.NewTodoFilterHistoryService(repo, logger)

		result, err := service.FetchLatestFilters(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedHistories, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := new(testutils.MockTodoFilterHistoryRepository)
		repo.On("FetchLatestFilters", mock.Anything, 5).Return([]*ent.TodoFilterHistory(nil), assert.AnError)

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := services.NewTodoFilterHistoryService(repo, logger)

		result, err := service.FetchLatestFilters(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestTodoFilterHistoryService_SaveFilterHistory(t *testing.T) {
	t.Run("リポジトリの SaveFilterHistory を呼び出し、結果を返すこと", func(t *testing.T) {
		repo := new(testutils.MockTodoFilterHistoryRepository)
		query := "test query"
		functionName := "testFunction"
		args := map[string]interface{}{"key": "value"}
		resultTodoIds := []int{1, 2, 3}
		expectedHistory := &ent.TodoFilterHistory{
			Query:        query,
			FunctionName: functionName,
			Args:         args,
		}

		repo.On("SaveFilterHistory", mock.Anything, query, &functionName, args, resultTodoIds).Return(expectedHistory, nil)

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := services.NewTodoFilterHistoryService(repo, logger)

		result, err := service.SaveFilterHistory(context.Background(), query, &functionName, args, resultTodoIds)

		assert.NoError(t, err)
		assert.Equal(t, expectedHistory, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリがエラーを返した場合、そのままエラーを返すこと", func(t *testing.T) {
		repo := new(testutils.MockTodoFilterHistoryRepository)
		repo.On("SaveFilterHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return((*ent.TodoFilterHistory)(nil), assert.AnError)

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := services.NewTodoFilterHistoryService(repo, logger)

		result, err := service.SaveFilterHistory(context.Background(), "", nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}
