package services

import (
	"todo-app/dto"
	"todo-app/ent"
)

type TodoService struct{}

func NewTodoService() *TodoService {
	return &TodoService{}
}

// EntityToDTO は単一の ent.Todo を dto.ListTodoResponse に変換します
func (s *TodoService) EntityToDTO(t *ent.Todo) dto.ListTodoResponse {
	return dto.ListTodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// EntitiesToDTOs はスライスを一括変換します
func (s *TodoService) EntitiesToDTOs(todos []*ent.Todo) []dto.ListTodoResponse {
	res := make([]dto.ListTodoResponse, len(todos))
	for i, t := range todos {
		res[i] = s.EntityToDTO(t)
	}
	return res
}
