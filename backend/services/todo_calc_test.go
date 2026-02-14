package services_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"todo-app/services"

	"github.com/stretchr/testify/assert"
)

func TestTodoService_CalculatePagination(t *testing.T) {
	t.Run("ページネーション情報が正しく計算されること", func(t *testing.T) {
		tests := []struct {
			name        string
			count       int
			currentPage int
			limit       int
			expected    struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}
		}{
			{
				name:        "1ページ目、次あり",
				count:       50,
				currentPage: 1,
				limit:       20,
				expected: struct {
					totalPages int
					hasNext    bool
					hasPrev    bool
				}{3, true, false},
			},
			{
				name:        "2ページ目、前後あり",
				count:       50,
				currentPage: 2,
				limit:       20,
				expected: struct {
					totalPages int
					hasNext    bool
					hasPrev    bool
				}{3, true, true},
			},
			{
				name:        "最終ページ、次なし",
				count:       50,
				currentPage: 3,
				limit:       20,
				expected: struct {
					totalPages int
					hasNext    bool
					hasPrev    bool
				}{3, false, true},
			},
			{
				name:        "データなし",
				count:       0,
				currentPage: 1,
				limit:       20,
				expected: struct {
					totalPages int
					hasNext    bool
					hasPrev    bool
				}{1, false, false},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := &spyTodoRepo{
					count: func(includeDone bool) (int, error) {
						return tt.count, nil
					},
				}
				ctx := context.Background()
				service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

				pagination, err := service.CalculatePagination(ctx, tt.currentPage, tt.limit, false)

				assert.NoError(t, err)
				assert.Equal(t, tt.expected.totalPages, pagination.TotalPages)
				assert.Equal(t, tt.expected.hasNext, pagination.HasNext)
				assert.Equal(t, tt.expected.hasPrev, pagination.HasPrev)
				assert.Equal(t, tt.currentPage, pagination.CurrentPage)
				assert.Equal(t, tt.limit, pagination.Limit)
			})
		}
	})

	t.Run("カウント取得でエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		repo := &spyTodoRepo{
			count: func(includeDone bool) (int, error) {
				return 0, errors.New("db error")
			},
		}
		ctx := context.Background()
		service := services.NewTodoService(slog.New(slog.NewTextHandler(io.Discard, nil)), repo)

		pagination, err := service.CalculatePagination(ctx, 1, 20, false)

		assert.Error(t, err)
		assert.Nil(t, pagination)
		assert.Equal(t, "db error", err.Error())
	})
}
