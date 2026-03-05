package repositories

import (
	"context"
	"todo-app/ent"
	"todo-app/ent/todofilterhistory"
	"todo-app/ent/user"
)

type ITodoFilterHistoryRepository interface {
	FetchLatestFilters(ctx context.Context, limit int) ([]*ent.TodoFilterHistory, error)
}

type TodoFilterHistoryRepository struct {
	base *BaseRepository
}

func NewTodoFilterHistoryRepository(client *ent.Client) *TodoFilterHistoryRepository {
	return &TodoFilterHistoryRepository{
		base: NewBaseRepository(client),
	}
}

func (r *TodoFilterHistoryRepository) FetchLatestFilters(ctx context.Context, limit int) ([]*ent.TodoFilterHistory, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)
	return client.TodoFilterHistory.Query().
		Where(todofilterhistory.HasUserWith(user.ID(u.ID))).
		Order(ent.Desc(todofilterhistory.FieldCreatedAt)).
		Limit(limit).
		All(ctx)
}
