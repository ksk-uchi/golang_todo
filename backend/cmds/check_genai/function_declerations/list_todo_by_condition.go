package check_genai_function_declerations

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type SearchTodosArgs struct {
	DoneFrom string `json:"done_from" description:"完了日時の開始範囲 (RFC3339形式)"`
	DoneTo   string `json:"done_to" description:"完了日時の終了範囲 (RFC3339形式)"`
}

func ListTodosByCondition(ctx context.Context, args SearchTodosArgs) {
	fmt.Println("This function is called.")
	fmt.Println(args)
}

var ListTodosByConditionDeclaration = &genai.FunctionDeclaration{
	Name:        "listTodosByCondition",
	Description: "タスクの完了日に対する範囲指定を行うための `done_from` と `done_to` を生成する関数",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"done_from": {
				Type: genai.TypeString,
			},
			"done_to": {
				Type: genai.TypeString,
			},
		},
	},
}
