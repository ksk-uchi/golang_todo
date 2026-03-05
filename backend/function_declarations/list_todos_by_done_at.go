package function_declerations

import (
	"context"
	"time"
	"todo-app/ent"
	"todo-app/repositories"

	"google.golang.org/genai"
)

type ListTodosByDoneAtArgs struct {
	DoneFrom string `json:"done_from"`
	DoneTo   string `json:"done_to"`
}

func ListTodosByDoneAt(ctx context.Context, repo repositories.ITodoRepository, args ListTodosByDoneAtArgs) ([]*ent.Todo, error) {
	var doneFrom, doneTo *time.Time

	if args.DoneFrom != "" {
		t, err := time.Parse(time.RFC3339, args.DoneFrom)
		if err != nil {
			return nil, err
		}
		doneFrom = &t
	}

	if args.DoneTo != "" {
		t, err := time.Parse(time.RFC3339, args.DoneTo)
		if err != nil {
			return nil, err
		}
		doneTo = &t
	}

	return repo.FetchTodosByDoneAt(ctx, doneFrom, doneTo)
}

var ListTodosByDoneAtDeclaration = &genai.FunctionDeclaration{
	Name:        "ListTodosByDoneAt",
	Description: "ToDo の完了日時に対して範囲指定で ToDo 検索する",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"done_from": {
				Type:        genai.TypeString,
				Description: "完了日時の開始範囲 (RFC3339形式)",
			},
			"done_to": {
				Type:        genai.TypeString,
				Description: "完了日時の終了範囲 (RFC3339形式)",
			},
		},
	},
}
