package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type GetimgAI struct {
	Cookie string
}

type GetimgAIGenerateImageRequest struct {
	Tool      string `json:"tool"`
	Style     string `json:"style"`
	Prompt    string `json:"prompt"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	NumImages int    `json:"num_images"`
}

type GetimgAIImageItem struct {
	JpegURL string `json:"jpegUrl"`
}

type GetimgAIGenerateImageResponse struct {
	Images []GetimgAIImageItem `json:"images"`
	// Ignore other fields
}

func NewGetimgAI(cookie string) *GetimgAI {
	return &GetimgAI{Cookie: cookie}
}

func (g *GetimgAI) CreateImage(req *openai.ImageRequest) (*openai.ImageResponse, error) {
	url := "https://getimg.ai/api/pipelines/" + req.Model
	width, height := 1024, 1024
	if req.Size != "" {
		parts := strings.Split(req.Size, "x")
		if len(parts) == 2 {
			width, _ = strconv.Atoi(parts[0])
			height, _ = strconv.Atoi(parts[1])
		}
	}

	if req.N == 0 {
		req.N = 1
	}
	// Create the request body
	payload := GetimgAIGenerateImageRequest{
		Tool:      "generator",
		Style:     "photorealism",
		Prompt:    req.Prompt,
		Width:     width,
		Height:    height,
		NumImages: req.N,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req2, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Referer", "https://getimg.ai/text-to-image")
	req2.Header.Set("Origin", "https://getimg.ai")
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req2.Header.Set("Cookie", g.Cookie)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var generateResp []GetimgAIGenerateImageResponse
	err = json.Unmarshal(respBody, &generateResp)
	if err != nil {

		return nil, err
	}

	images := make([]openai.ImageResponseDataInner, 0, len(generateResp[0].Images))
	for _, img := range generateResp[0].Images {
		images = append(images, openai.ImageResponseDataInner{URL: img.JpegURL})
	}
	resp2 := openai.ImageResponse{
		Created: time.Now().Unix(),
		Data:    images,
	}
	return &resp2, nil
}

var _ OpenAI = &GetimgAI{}
