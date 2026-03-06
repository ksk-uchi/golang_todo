package utils

import (
	"context"
	"fmt"
	"os"
	"sync"
	"todo-app/services"

	"google.golang.org/genai"
)

type IAIFactory interface {
	GetGeminiClient(ctx context.Context) (services.IGenAIClient, error)
}

type AIFactory struct {
	geminiClient services.IGenAIClient
	geminiOnce   sync.Once
	geminiErr    error
}

func NewAIFactory() IAIFactory {
	return &AIFactory{}
}

// genAIClientWrapper is moved from services package
type genAIClientWrapper struct {
	client *genai.Client
}

// GenerateContent implements services.IGenAIClient
func (w *genAIClientWrapper) GenerateContent(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	return w.client.Models.GenerateContent(ctx, model, contents, config)
}

func (f *AIFactory) GetGeminiClient(ctx context.Context) (services.IGenAIClient, error) {
	f.geminiOnce.Do(func() {
		apiKey := os.Getenv("GOOGLE_API_KEY")
		if apiKey == "" {
			f.geminiErr = fmt.Errorf("GOOGLE_API_KEY is not set")
			return
		}
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			f.geminiErr = err
			return
		}
		f.geminiClient = &genAIClientWrapper{client: client}
	})
	return f.geminiClient, f.geminiErr
}
