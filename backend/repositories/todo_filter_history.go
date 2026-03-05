package repositories

import (
	"context"
	"todo-app/ent"
	"todo-app/ent/todofilterhistory"
	"todo-app/ent/user"
)

type ITodoFilterHistoryRepository interface {
	FetchLatestFilters(ctx context.Context, userID int, limit int) ([]*ent.TodoFilterHistory, error)
}

type TodoFilterHistoryRepository struct {
	client *ent.Client
}

func NewTodoFilterHistoryRepository(client *ent.Client) *TodoFilterHistoryRepository {
	return &TodoFilterHistoryRepository{client: client}
}

func (r *TodoFilterHistoryRepository) FetchLatestFilters(ctx context.Context, userID int, limit int) ([]*ent.TodoFilterHistory, error) {
	return r.client.TodoFilterHistory.Query().
		Where(todofilterhistory.HasUserWith(user.ID(userID))).
		Order(ent.Desc(todofilterhistory.FieldCreatedAt)).
		Limit(limit).
		All(ctx)
}
