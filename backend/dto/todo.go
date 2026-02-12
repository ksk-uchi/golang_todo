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

type ListTodoResponseDto struct {
	Data []TodoDto `json:"data"`
}

func EntitiesToTodoDtoSlice(todos []*ent.Todo) ListTodoResponseDto {
	todoDtos := make([]TodoDto, len(todos))
	for i, t := range todos {
		todoDtos[i] = EntityToTodoDto(t)
	}
	return ListTodoResponseDto{
		Data: todoDtos,
	}
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
