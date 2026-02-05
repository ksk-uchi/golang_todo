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

func (r *TodoRepository) FindTodo(id int) (*ent.Todo, error) {
	return r.client.Todo.Get(r.ctx, id)
}

func (r *TodoRepository) CreateTodo(title string, description string) (*ent.Todo, error) {
	return r.client.Todo.Create().
		SetTitle(title).
		SetDescription(description).
		Save(r.ctx)
}

func (r *TodoRepository) UpdateTodo(id int, title *string, description *string) (*ent.Todo, error) {
	return r.client.Todo.UpdateOneID(id).
		SetNillableTitle(title).
		SetNillableDescription(description).
		Save(r.ctx)
}

func (r *TodoRepository) DeleteTodo(id int) error {
	return r.client.Todo.DeleteOneID(id).Exec(r.ctx)
}
