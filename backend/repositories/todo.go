package repositories

import (
	"context"
	"log/slog"
	"time"
	"todo-app/ent"
	"todo-app/ent/todo"
	"todo-app/ent/user"
)

type ITodoRepository interface {
	FetchTodos(ctx context.Context, limit int, offset int, includeDone bool) ([]*ent.Todo, error)
	GetTodoCount(ctx context.Context, includeDone bool) (int, error)
	FindTodo(ctx context.Context, id int) (*ent.Todo, error)
	GetTodoForUpdate(ctx context.Context, id int) (*ent.Todo, error)
	CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error)
	UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error)
	UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error)
	DeleteTodo(ctx context.Context, id int) error
	FetchTodosByDoneAt(ctx context.Context, doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error)
}

type TodoRepository struct {
	base   *BaseRepository
	logger *slog.Logger
}

func NewTodoRepository(client *ent.Client, logger *slog.Logger) *TodoRepository {
	return &TodoRepository{
		base:   NewBaseRepository(client),
		logger: logger,
	}
}

func (r *TodoRepository) FetchTodos(ctx context.Context, limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)
	query := client.Todo.Query().
		Where(todo.HasUserWith(user.ID(u.ID)))

	if !includeDone {
		query.Where(todo.DoneAtIsNil())
	}

	return query.
		Order(ent.Desc(todo.FieldUpdatedAt), ent.Desc(todo.FieldID)).
		Limit(limit).
		Offset(offset).
		All(ctx)
}

func (r *TodoRepository) GetTodoCount(ctx context.Context, includeDone bool) (int, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return 0, err
	}
	client := r.base.getClient(ctx)
	query := client.Todo.Query().
		Where(todo.HasUserWith(user.ID(u.ID)))

	if !includeDone {
		query.Where(todo.DoneAtIsNil())
	}

	return query.Count(ctx)
}

func (r *TodoRepository) FindTodo(ctx context.Context, id int) (*ent.Todo, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)
	return client.Todo.Query().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		Only(ctx)
}

func (r *TodoRepository) CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)
	return client.Todo.Create().
		SetTitle(title).
		SetDescription(description).
		SetUser(u).
		Save(ctx)
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)

	return client.Todo.UpdateOneID(id).
		Where(todo.HasUserWith(user.ID(u.ID))).
		SetNillableTitle(title).
		SetNillableDescription(description).
		Save(ctx)
}

func (r *TodoRepository) UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)

	update := client.Todo.UpdateOneID(id).
		Where(todo.HasUserWith(user.ID(u.ID)))

	if isDone {
		update.SetDoneAt(time.Now())
	} else {
		update.ClearDoneAt()
	}

	return update.Save(ctx)
}

func (r *TodoRepository) GetTodoForUpdate(ctx context.Context, id int) (*ent.Todo, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)

	query := client.Todo.Query().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		ForUpdate()

	res, err := query.Only(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *TodoRepository) FetchTodosByDoneAt(ctx context.Context, doneFrom *time.Time, doneTo *time.Time) ([]*ent.Todo, error) {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := r.base.getClient(ctx)
	query := client.Todo.Query().
		Where(todo.HasUserWith(user.ID(u.ID)))

	if doneFrom != nil {
		query.Where(todo.DoneAtGTE(*doneFrom))
	}
	if doneTo != nil {
		query.Where(todo.DoneAtLTE(*doneTo))
	}

	return query.
		Order(ent.Desc(todo.FieldDoneAt), ent.Desc(todo.FieldID)).
		All(ctx)
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id int) error {
	u, err := r.base.getUser(ctx)
	if err != nil {
		return err
	}
	client := r.base.getClient(ctx)
	// Verify ownership and delete
	n, err := client.Todo.Delete().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}
