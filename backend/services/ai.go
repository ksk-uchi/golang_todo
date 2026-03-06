package services

import (
	"context"
	"fmt"
	"time"
	"todo-app/dto"
	"todo-app/ent"
	function_declerations "todo-app/function_declarations"
	"todo-app/repositories"
	"todo-app/utils"

	"google.golang.org/genai"
)

type AIService struct {
	repo repositories.ITodoRepository
}

func NewAIService(repo repositories.ITodoRepository) *AIService {
	return &AIService{repo: repo}
}

func (s *AIService) DecideFilterTodosFunction(ctx context.Context, aiClient utils.IGenAIClient, query string) (*dto.AIFilterDto, error) {
	parts := []*genai.Part{
		{Text: time.Now().Format("現在2006年1月2日15:04:05です。")},
		{Text: "使用できる Tool が無い場合は「対応できる Tool がありません」とだけ回答するようにしてください。"},
		{Text: query},
	}

	result, err := aiClient.GenerateContent(ctx,
		"gemini-3-flash-preview",
		[]*genai.Content{{Parts: parts}},
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						function_declerations.ListTodosByDoneAtDeclaration,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, nil
	}

	fc := result.Candidates[0].Content.Parts[0].FunctionCall
	if fc == nil {
		return nil, nil
	}

	return &dto.AIFilterDto{
		FunctionName: fc.Name,
		Args:         fc.Args,
	}, nil
}

func (s *AIService) FilterTodos(ctx context.Context, functionName string, args interface{}) ([]*ent.Todo, error) {
	if functionName == "ListTodosByDoneAt" {
		m, ok := args.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid args type")
		}

		listArgs := function_declerations.ListTodosByDoneAtArgs{}
		if v, ok := m["done_from"].(string); ok {
			listArgs.DoneFrom = v
		}
		if v, ok := m["done_to"].(string); ok {
			listArgs.DoneTo = v
		}

		return function_declerations.ListTodosByDoneAt(ctx, s.repo, listArgs)
	}

	return nil, fmt.Errorf("unknown function: %s", functionName)
}
