package services

import (
	"todo-app/ent"
	"todo-app/ent/todo"
	"todo-app/middlewares"

	"github.com/labstack/echo/v5"
)

func NewTodoService() *TodoService {
	return &TodoService{}
}

type TodoService struct{}

func (s *TodoService) GetTodoSlice(c *echo.Context) ([]*ent.Todo, error) {
	tx := middlewares.GetTx(c)

	ctx := c.Request().Context()
	return tx.Todo.Query().
		Order(ent.Desc(todo.FieldCreatedAt)).
		All(ctx)
}
