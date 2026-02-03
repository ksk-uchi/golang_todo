package services

import (
	"todo-app/dto"
	"todo-app/ent"
)

type ITodoService interface {
	EntitiesToDTOs(todos []*ent.Todo) []dto.ListTodoResponse
}

type TodoService struct {
	client *ent.Client
}

func NewTodoService(client *ent.Client) *TodoService { return &TodoService{client: client} }

func (s *TodoService) EntitiesToDTOs(todos []*ent.Todo) []dto.ListTodoResponse {
	res := make([]dto.ListTodoResponse, len(todos))
	for i, t := range todos {
		res[i] = dto.ListTodoResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}
	return res
}
