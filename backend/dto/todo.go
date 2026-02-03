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
		res[i] = TodoDto{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}
	return res
}
