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

func NewOpenAITranslator(baseURL, apiKey, model, promptTemplate string) *OpenAITranslator {
	return &OpenAITranslator{
		BaseURL:        baseURL,
		APIKey:         apiKey,
		Model:          model,
		PromptTemplate: promptTemplate,
	}
}

func (t *OpenAITranslator) Translate(content string) (string, error) {
	config := openai.DefaultConfig(t.APIKey)
	config.BaseURL = t.BaseURL
	c := openai.NewClientWithConfig(config)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: t.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(t.PromptTemplate, content),
			},
		},
	}
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response")
}

var _ Translator = &OpenAITranslator{}
