package repositories

import (
	"context"
	"todo-app/ent"
	"todo-app/ent/user"
)

type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*ent.User, error)
}

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
}
