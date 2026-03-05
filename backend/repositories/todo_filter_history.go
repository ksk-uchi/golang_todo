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
	client *ent.Client
}

func NewTodoFilterHistoryRepository(client *ent.Client) *TodoFilterHistoryRepository {
	return &TodoFilterHistoryRepository{client: client}
}

func (r *TodoFilterHistoryRepository) getClient(ctx context.Context) *ent.Client {
	tx := ent.TxFromContext(ctx)
	if tx != nil {
		return tx.Client()
	}
	return r.client
}

func (r *TodoFilterHistoryRepository) getUser(ctx context.Context) (*ent.User, error) {
	u := ctx.Value("user")
	if u == nil {
		return nil, &ent.NotFoundError{}
	}
	return u.(*ent.User), nil
}

func (r *TodoFilterHistoryRepository) FetchLatestFilters(ctx context.Context, limit int) ([]*ent.TodoFilterHistory, error) {
	u, err := r.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.getClient(ctx)
	return client.TodoFilterHistory.Query().
		Where(todofilterhistory.HasUserWith(user.ID(u.ID))).
		Order(ent.Desc(todofilterhistory.FieldCreatedAt)).
		Limit(limit).
		All(ctx)
}
