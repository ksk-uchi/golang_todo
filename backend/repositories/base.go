package repositories

import (
	"context"
	"todo-app/ent"
)

type BaseRepository struct {
	client *ent.Client
}

func NewBaseRepository(client *ent.Client) *BaseRepository {
	return &BaseRepository{client: client}
}

func (r *BaseRepository) getClient(ctx context.Context) *ent.Client {
	tx := ent.TxFromContext(ctx)
	if tx != nil {
		return tx.Client()
	}
	return r.client
}

func (r *BaseRepository) getUser(ctx context.Context) (*ent.User, error) {
	u := ctx.Value("user")
	if u == nil {
		return nil, &ent.NotFoundError{}
	}
	return u.(*ent.User), nil
}
