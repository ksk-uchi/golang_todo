package repositories

import (
	"context"
	"todo-app/ent"
	"todo-app/ent/todo"
	"todo-app/ent/user"
)

type TodoRepository struct {
	ctx    context.Context
	client *ent.Client
}

func NewTodoRepository(ctx context.Context, client *ent.Client) *TodoRepository {
	return &TodoRepository{ctx: ctx, client: client}
}

func (r *TodoRepository) getUser() (*ent.User, error) {
	u := r.ctx.Value("user")
	if u == nil {
		return nil, &ent.NotFoundError{}
	}
	return u.(*ent.User), nil
}

func (r *TodoRepository) FetchTodos(limit int, offset int) ([]*ent.Todo, error) {
	u, err := r.getUser()
	if err != nil {
		return nil, err
	}
	return r.client.Todo.Query().
		Where(todo.HasUserWith(user.ID(u.ID))).
		Order(ent.Desc(todo.FieldUpdatedAt), ent.Desc(todo.FieldID)).
		Limit(limit).
		Offset(offset).
		All(r.ctx)
}

func (r *TodoRepository) GetTodoCount() (int, error) {
	u, err := r.getUser()
	if err != nil {
		return 0, err
	}
	return r.client.Todo.Query().
		Where(todo.HasUserWith(user.ID(u.ID))).
		Count(r.ctx)
}

func (r *TodoRepository) FindTodo(id int) (*ent.Todo, error) {
	u, err := r.getUser()
	if err != nil {
		return nil, err
	}
	return r.client.Todo.Query().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		Only(r.ctx)
}

func (r *TodoRepository) CreateTodo(title string, description string) (*ent.Todo, error) {
	u, err := r.getUser()
	if err != nil {
		return nil, err
	}
	return r.client.Todo.Create().
		SetTitle(title).
		SetDescription(description).
		SetUser(u).
		Save(r.ctx)
}

func (r *TodoRepository) UpdateTodo(id int, title *string, description *string) (*ent.Todo, error) {
	u, err := r.getUser()
	if err != nil {
		return nil, err
	}
	// Verify ownership before update
	exists, err := r.client.Todo.Query().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		Exist(r.ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &ent.NotFoundError{}
	}

	return r.client.Todo.UpdateOneID(id).
		Where(todo.HasUserWith(user.ID(u.ID))).
		SetNillableTitle(title).
		SetNillableDescription(description).
		Save(r.ctx)
}

func (r *TodoRepository) DeleteTodo(id int) error {
	u, err := r.getUser()
	if err != nil {
		return err
	}
	// Verify ownership and delete
	n, err := r.client.Todo.Delete().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		Exec(r.ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}
