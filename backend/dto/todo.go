package dto

import (
	"time"
	"todo-app/ent"
)

type TodoDto struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DoneAt      *time.Time `json:"done_at"`
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

type AIFilterDto struct {
	FunctionName string                 `json:"function_name"`
	Args         map[string]interface{} `json:"args"`
}

type TodoFilterHistoryQueryDto struct {
	Query string `json:"query"`
	ID    string `json:"id"`
}

type ListTodoFilterHistoriesResponseDto struct {
	Queries []TodoFilterHistoryQueryDto `json:"queries"`
}

func EntityToTodoDto(todo *ent.Todo) TodoDto {
	return TodoDto{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
		DoneAt:      todo.DoneAt,
	}
}

func EntitiesToTodoFilterHistoryQueryDtos(histories []*ent.TodoFilterHistory) []TodoFilterHistoryQueryDto {
	dtos := make([]TodoFilterHistoryQueryDto, len(histories))
	for i, h := range histories {
		dtos[i] = TodoFilterHistoryQueryDto{
			Query: h.Query,
			ID:    h.ID.String(),
		}
	}
	return dtos
}
