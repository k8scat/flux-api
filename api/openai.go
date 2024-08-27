package api

// Refer to https://platform.openai.com/docs/api-reference/images/create

type GenerateImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageItem struct {
	URL string `json:"url"`
}

type GenerateImageResponse struct {
	Created int64       `json:"created"`
	Data    []ImageItem `json:"data"`
}

type OpenAI interface {
	GenerateImage(req *GenerateImageRequest) (*GenerateImageResponse, error)
}
