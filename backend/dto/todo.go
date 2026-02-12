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
	Data       []TodoDto      `json:"data"`
	Pagination *PaginationDto `json:"pagination"`
}

type PaginationDto struct {
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
	Limit       int  `json:"limit"`
}

func EntitiesToTodoDtoSlice(todos []*ent.Todo, pagination *PaginationDto) ListTodoResponseDto {
	todoDtos := make([]TodoDto, len(todos))
	for i, t := range todos {
		todoDtos[i] = EntityToTodoDto(t)
	}
	return ListTodoResponseDto{
		Data:       todoDtos,
		Pagination: pagination,
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
