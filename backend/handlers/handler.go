package handlers

import "todo-app/ent"

type Handler struct {
	Ent *ent.Client
}

func NewHandler(client *ent.Client) *Handler {
	return &Handler{Ent: client}
}
