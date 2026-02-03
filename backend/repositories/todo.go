package repositories

import (
	"context"
	"todo-app/ent"
	"todo-app/ent/todo"
)

type TodoRepository struct {
	client *ent.Client
}

func NewTodoRepository(client *ent.Client) *TodoRepository {
	return &TodoRepository{client: client}
}

func (r *TodoRepository) FetchAllTodo(ctx context.Context) ([]*ent.Todo, error) {
	return r.client.Todo.Query().Order(ent.Desc(todo.FieldCreatedAt)).All(ctx)
}
