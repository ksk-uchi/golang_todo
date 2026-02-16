package repositories

import (
	"context"
	"log/slog"
	"strings"
	"time"
	"todo-app/ent"
	"todo-app/ent/todo"
	"todo-app/ent/user"
)

type TodoRepository struct {
	logger *slog.Logger
}

func NewTodoRepository(logger *slog.Logger) *TodoRepository {
	return &TodoRepository{logger: logger}
}

func (r *TodoRepository) getUser(ctx context.Context) (*ent.User, error) {
	u := ctx.Value("user")
	if u == nil {
		return nil, &ent.NotFoundError{}
	}
	return u.(*ent.User), nil
}

// ... Check existing content ...

func (r *TodoRepository) FetchTodos(ctx context.Context, limit int, offset int, includeDone bool) ([]*ent.Todo, error) {
	u, err := r.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := getEntClient(ctx)
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
	u, err := r.getUser(ctx)
	if err != nil {
		return 0, err
	}
	client := getEntClient(ctx)
	query := client.Todo.Query().
		Where(todo.HasUserWith(user.ID(u.ID)))

	if !includeDone {
		query.Where(todo.DoneAtIsNil())
	}

	return query.Count(ctx)
}

func (r *TodoRepository) FindTodo(ctx context.Context, id int) (*ent.Todo, error) {
	u, err := r.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := getEntClient(ctx)
	return client.Todo.Query().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		Only(ctx)
}

func (r *TodoRepository) CreateTodo(ctx context.Context, title string, description string) (*ent.Todo, error) {
	u, err := r.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := getEntClient(ctx)
	return client.Todo.Create().
		SetTitle(title).
		SetDescription(description).
		SetUser(u).
		Save(ctx)
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, id int, title *string, description *string) (*ent.Todo, error) {
	u, err := r.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := getEntClient(ctx)

	return client.Todo.UpdateOneID(id).
		Where(todo.HasUserWith(user.ID(u.ID))).
		SetNillableTitle(title).
		SetNillableDescription(description).
		Save(ctx)
}

func (r *TodoRepository) UpdateDoneStatus(ctx context.Context, id int, isDone bool) (*ent.Todo, error) {
	u, err := r.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := getEntClient(ctx)

	// Helper function to execute update
	executeUpdate := func(t *ent.Todo) (*ent.Todo, error) {
		if isDone {
			if t.DoneAt != nil {
				return t, nil
			}
			return t.Update().SetDoneAt(time.Now()).Save(ctx)
		} else {
			if t.DoneAt == nil {
				return t, nil
			}
			return t.Update().ClearDoneAt().Save(ctx)
		}
	}

	// Try with ForUpdate
	t, err := client.Todo.Query().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		ForUpdate().
		Only(ctx)

	if err != nil {
		// Fallback for SQLite which doesn't support ForUpdate
		if strings.Contains(err.Error(), "not supported") {
			t, err = client.Todo.Query().
				Where(todo.ID(id)).
				Where(todo.HasUserWith(user.ID(u.ID))).
				Only(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return executeUpdate(t)
}

func (r *TodoRepository) GetTodoForUpdate(ctx context.Context, id int) (*ent.Todo, error) {
	u, err := r.getUser(ctx)
	if err != nil {
		return nil, err
	}
	client := getEntClient(ctx)

	query := client.Todo.Query().
		Where(todo.ID(id)).
		Where(todo.HasUserWith(user.ID(u.ID))).
		ForUpdate()

	res, err := query.Only(ctx)
	if err != nil {
		// Fallback for SQLite
		if strings.Contains(err.Error(), "not supported") {
			return client.Todo.Query().
				Where(todo.ID(id)).
				Where(todo.HasUserWith(user.ID(u.ID))).
				Only(ctx)
		}
		return nil, err
	}
	return res, nil
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id int) error {
	u, err := r.getUser(ctx)
	if err != nil {
		return err
	}
	client := getEntClient(ctx)
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
