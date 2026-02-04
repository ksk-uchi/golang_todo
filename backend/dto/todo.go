package dto

import (
	"time"
	"todo-app/ent"
)

type TodoDto struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func EntitiesToTodoDtoSlice(todos []*ent.Todo) []TodoDto {
	res := make([]TodoDto, len(todos))
	for i, t := range todos {
		res[i] = EntityToTodoDto(t)
	}
	return res
}

type CreateTodoDto struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=200"`
}

func EntityToTodoDto(todo *ent.Todo) TodoDto {
	return TodoDto{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}
