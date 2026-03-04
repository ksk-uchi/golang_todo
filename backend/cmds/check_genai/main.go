package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"todo-app/cmds/check_genai/function_declerations"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	_ "github.com/go-sql-driver/mysql"
)

func processUserInput(ctx context.Context, client *genai.Client, inputText string) error {
	// models, err := client.Models.List(ctx, &genai.ListModelsConfig{})
	// if err != nil {
	// 	return err
	// }
	// for _, model := range models.Items {
	// 	fmt.Println(model.Name)
	// }
	parts := []*genai.Part{
		{Text: time.Now().Format("現在2006年1月2日15:04:05です。")},
		{Text: inputText},
	}
	result, err := client.Models.GenerateContent(ctx,
		"gemini-3-flash-preview",
		[]*genai.Content{{Parts: parts}},
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						function_declerations.ListTodosByConditionDeclaration,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("no response from gemini")
	}
	fc := result.Candidates[0].Content.Parts[0].FunctionCall
	if fc == nil {
		return fmt.Errorf("function call not found")
	}
	if fc.Name == "listTodosByCondition" {
		function_declerations.ListTodosByCondition(ctx, function_declerations.SearchTodosArgs{
			DoneFrom: fc.Args["done_from"].(string),
			DoneTo:   fc.Args["done_to"].(string),
		})
	}
	return nil
}

func main() {
	if err := godotenv.Load("envs/local.env"); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GOOGLE_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	text := "直近完了したタスクを見たい"
	// text := "今日はなんの記念日ですか？"
	if err := processUserInput(ctx, client, text); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
