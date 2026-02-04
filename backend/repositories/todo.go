package repositories

import (
	"context"
	"todo-app/ent"
	"todo-app/ent/todo"
)

type TodoRepository struct {
	ctx    context.Context
	client *ent.Client
}

func NewTodoRepository(ctx context.Context, client *ent.Client) *TodoRepository {
	return &TodoRepository{ctx: ctx, client: client}
}

func (r *TodoRepository) FetchAllTodo() ([]*ent.Todo, error) {
	return r.client.Todo.Query().Order(ent.Desc(todo.FieldCreatedAt)).All(r.ctx)
}
