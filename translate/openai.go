package translate

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type OpenAITranslator struct {
	BaseURL        string
	APIKey         string
	Model          string
	PromptTemplate string
}

func NewOpenAITranslator(baseURL, apiKey, model string) *OpenAITranslator {
	return &OpenAITranslator{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
	}
}

func (t *OpenAITranslator) Translate(content string) (string, error) {
	config := openai.DefaultConfig(t.APIKey)
	config.BaseURL = t.BaseURL
	c := openai.NewClientWithConfig(config)
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:  t.Model,
		Prompt: fmt.Sprintf(t.PromptTemplate, content),
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}

var _ Translator = &OpenAITranslator{}
