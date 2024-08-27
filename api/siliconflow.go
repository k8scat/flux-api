package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Refer to https://docs.siliconflow.cn/reference/black-forest-labsflux1-schnell
type SiliconFlow struct {
	APIKey string
}

type SiliconFlowGenerateImageRequest struct {
	Prompt            string `json:"prompt"`
	ImageSize         string `json:"image_size"`
	Seed              int64  `json:"seed"`
	NumInferenceSteps int    `json:"num_inference_steps"`
}

type SiliconFlowGenerateImageResponse struct {
	Images []ImageItem `json:"images"`
	// Ignore other fields
}

func NewSiliconFlow(apiKey string) *SiliconFlow {
	return &SiliconFlow{APIKey: apiKey}
}

func (s *SiliconFlow) GenerateImage(req *GenerateImageRequest) (*GenerateImageResponse, error) {
	url := "https://api.siliconflow.cn/v1/black-forest-labs/" + req.Model + "/text-to-image"
	req2 := SiliconFlowGenerateImageRequest{
		Prompt:            req.Prompt,
		ImageSize:         req.Size,
		Seed:              rand.Int63n(9999999999),
		NumInferenceSteps: 20,
	}
	if req2.ImageSize == "" {
		req2.ImageSize = "1024x1024"
	}

	b, _ := json.Marshal(req2)
	// Make a POST request to the URL with req2
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert the response body to a string
	var imgResp SiliconFlowGenerateImageResponse
	if err := json.Unmarshal(body, &imgResp); err != nil {
		return nil, err
	}
	openaiResp := GenerateImageResponse{
		Created: time.Now().Unix(),
		Data:    imgResp.Images,
	}
	return &openaiResp, nil
}

var _ OpenAI = &SiliconFlow{}
