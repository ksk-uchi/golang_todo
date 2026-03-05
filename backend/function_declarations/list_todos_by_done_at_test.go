package function_declerations

import (
	"context"
	"testing"
	"time"
	"todo-app/ent"

	"github.com/stretchr/testify/assert"
)

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

func TestListTodosByDoneAt(t *testing.T) {
	repo := &spyTodoRepo{}
	ctx := context.Background()

	t.Run("success - full range", func(t *testing.T) {
		args := ListTodosByDoneAtArgs{
			DoneFrom: "2023-01-01T00:00:00Z",
			DoneTo:   "2023-01-31T23:59:59Z",
		}
		repo.fetchTodosByDoneAt = func(doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error) {
			assert.Equal(t, "2023-01-01T00:00:00Z", doneFrom.Format(time.RFC3339))
			assert.Equal(t, "2023-01-31T23:59:59Z", doneTo.Format(time.RFC3339))
			return []*ent.Todo{{ID: 1}}, nil
		}

		res, err := ListTodosByDoneAt(ctx, repo, args)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("success - empty range", func(t *testing.T) {
		args := ListTodosByDoneAtArgs{}
		repo.fetchTodosByDoneAt = func(doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error) {
			assert.Nil(t, doneFrom)
			assert.Nil(t, doneTo)
			return []*ent.Todo{}, nil
		}

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
