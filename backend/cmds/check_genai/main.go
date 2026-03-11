package main

import (
	"context"
	"fmt"
	"os"
	"time"
	function_declerations "todo-app/cmds/check_genai/function_declerations"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	_ "github.com/go-sql-driver/mysql"
)

func processUserInput(ctx context.Context, client *genai.Client, inputText string) error {
	parts := []*genai.Part{
		genai.NewPartFromText(time.Now().Format("現在2006年1月2日15:04:05です。")),
		genai.NewPartFromText("使用できる Tool が無い場合は「対応できる Tool がありません」とだけ回答するようにしてください。"),
		genai.NewPartFromText(inputText),
	}

	chat, err := client.Chats.Create(ctx, "gemini-3-flash-preview", &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: function_declerations.Filters,
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	result, err := chat.Send(ctx, parts...)
	if err != nil {
		return err
	}

	var currentIDs []int

	for {
		if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
			return fmt.Errorf("no response from gemini")
		}

		var functionResponseParts []*genai.Part
		hasFunctionCall := false

		// Print any text part or handle function calls
		for _, part := range result.Candidates[0].Content.Parts {
			if part.Text != "" {
				fmt.Println(part.Text)
			}
			if part.FunctionCall != nil {
				hasFunctionCall = true
				fc := part.FunctionCall
				var response map[string]any

				fmt.Println("FunctionCall: ", fc.Name)
				fmt.Println("Args: ", fc.Args)

				if fc.Name == "get_all_todo_ids" {
					currentIDs = function_declerations.GetAllTodoIDs()
				} else {
					// 抽出されたツール（filter_for_...等）が呼ばれたら絞り込み処理を実行
					currentIDs = function_declerations.RandomExcludeIDs(currentIDs)
				}
				fmt.Println("Response: ", currentIDs)

				response = map[string]any{"result": currentIDs}

				respPart := genai.NewPartFromFunctionResponse(fc.Name, response)
				if fc.ID != "" {
					respPart.FunctionResponse.ID = fc.ID
				}
				functionResponseParts = append(functionResponseParts, respPart)
			}
		}

		if !hasFunctionCall {
			fmt.Printf("最終的な数値の配列: %v\n", currentIDs)
			break
		}

		// Send the function execution results back to the model
		result, err = chat.Send(ctx, functionResponseParts...)
		if err != nil {
			return err
		}
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

	text := "今日は雨が降っています。男の子にはどんなタスクを割り振ったらよいでしょう？"
	if err := processUserInput(ctx, client, text); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
