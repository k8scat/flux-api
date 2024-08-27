package api

import "github.com/sashabaranov/go-openai"

// Refer to https://platform.openai.com/docs/api-reference/images/create

type OpenAI interface {
	CreateImage(req *openai.ImageRequest) (*openai.ImageResponse, error)
}
