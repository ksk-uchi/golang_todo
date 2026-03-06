package repositories

import (
	"context"
	"todo-app/ent"
	"todo-app/ent/todofilterhistory"
	"todo-app/ent/user"

	"github.com/google/uuid"
)

type ITodoFilterHistoryRepository interface {
	FetchLatestFilters(ctx context.Context, limit int) ([]*ent.TodoFilterHistory, error)
	SaveFilterHistory(ctx context.Context, query string, functionName *string, args map[string]interface{}, resultTodoIds []int) (*ent.TodoFilterHistory, error)
	GetFilterHistoryByQueryID(ctx context.Context, queryID uuid.UUID) (*ent.TodoFilterHistory, error)
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

func (r *TodoFilterHistoryRepository) SaveFilterHistory(ctx context.Context, query string, functionName *string, args map[string]interface{}, resultTodoIds []int) (*ent.TodoFilterHistory, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)
	create := client.TodoFilterHistory.Create().
		SetUserID(u.ID).
		SetQuery(query).
		SetNillableFunctionName(functionName).
		SetArgs(args).
		SetResultTodoIds(resultTodoIds)

	return create.Save(ctx)
}

func (r *TodoFilterHistoryRepository) GetFilterHistoryByQueryID(ctx context.Context, queryID uuid.UUID) (*ent.TodoFilterHistory, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)
	return client.TodoFilterHistory.Query().
		Where(todofilterhistory.IDEQ(queryID)).
		Where(todofilterhistory.HasUserWith(user.ID(u.ID))).
		Only(ctx)
}
