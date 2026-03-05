package function_declerations

import (
	"context"
	"testing"
	"time"
	"todo-app/ent"
	"todo-app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListTodosByDoneAt(t *testing.T) {
	repo := new(testutils.MockTodoRepository)
	ctx := context.Background()

	t.Run("success - full range", func(t *testing.T) {
		args := ListTodosByDoneAtArgs{
			DoneFrom: "2023-01-01T00:00:00Z",
			DoneTo:   "2023-01-31T23:59:59Z",
		}
		doneFrom, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
		doneTo, _ := time.Parse(time.RFC3339, "2023-01-31T23:59:59Z")
		repo.On("FetchTodosByDoneAt", mock.Anything, &doneFrom, &doneTo).Return([]*ent.Todo{{ID: 1}}, nil).Once()

		res, err := ListTodosByDoneAt(ctx, repo, args)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("success - empty range", func(t *testing.T) {
		args := ListTodosByDoneAtArgs{}
		repo.On("FetchTodosByDoneAt", mock.Anything, (*time.Time)(nil), (*time.Time)(nil)).Return([]*ent.Todo{}, nil).Once()

		res, err := ListTodosByDoneAt(ctx, repo, args)
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("error - invalid date", func(t *testing.T) {
		args := ListTodosByDoneAtArgs{
			DoneFrom: "invalid",
		}
		res, err := ListTodosByDoneAt(ctx, repo, args)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
