package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8scat/flux-api/api"
	"github.com/k8scat/flux-api/translate"
	"github.com/sashabaranov/go-openai"
)

func translatePrompt(prompt string) string {
	enable := strings.ToLower(os.Getenv("TRANSLATE_ENABLE"))
	if enable == "true" {
		translator := translate.NewOpenAITranslator(
			os.Getenv("TRANSLATE_API_BASE"),
			os.Getenv("TRANSLATE_API_KEY"),
			os.Getenv("TRANSLATE_MODEL"),
			os.Getenv("TRANSLATE_PROMPT_TEMPLATE"),
		)
		prompt, err := translator.Translate(prompt)
		if err != nil {
			log.Printf("failed to translate prompt: %v", err)
			return prompt
		}
		log.Printf("translated prompt: %s", prompt)
	}
	return prompt
}

func createImage(c *gin.Context, obj any) *openai.ImageResponse {
	if err := c.BindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}

	var req *openai.ImageRequest
	switch obj.(type) {
	case *openai.ImageRequest:
		req = obj.(*openai.ImageRequest)
	case *openai.ChatCompletionRequest:
		tmp := obj.(*openai.ChatCompletionRequest)
		if len(tmp.Messages) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return nil
		}
		req = &openai.ImageRequest{
			Model:  tmp.Model,
			Prompt: tmp.Messages[len(tmp.Messages)-1].Content,
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return nil
	}

	authHeader := c.GetHeader("Authorization")
	auth := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if auth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auth"})
		return nil
	}

	var client api.OpenAI
	if req.Model == "FLUX.1-schnell" {
		client = api.NewSiliconFlow(auth)
	} else if req.Model == "flux-v1" {
		client = api.NewGetimgAI(auth)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid model"})
		return nil
	}

	req.Prompt = translatePrompt(req.Prompt)
	resp, err := client.CreateImage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil
	}
	return resp
}

func main() {
	r := gin.Default()

	r.POST("/v1/images/generations", func(c *gin.Context) {
		var req openai.ImageRequest
		resp := createImage(c, &req)
		if resp != nil {
			c.JSON(http.StatusOK, resp)
		}
	})

	r.POST("/v1/chat/completions", func(c *gin.Context) {
		var req openai.ChatCompletionRequest
		resp := createImage(c, &req)
		if resp == nil {
			return
		}

		url := resp.Data[0].URL
		now := time.Now().Unix()
		ret := gin.H{
			"id":      fmt.Sprintf("chatcmpl-%d", now),
			"object":  "chat.completion",
			"created": now,
			"model":   req.Model,
			"choices": []gin.H{
				{
					"index": 0,
					"message": gin.H{
						"role":    "assistant",
						"content": fmt.Sprintf("![%s](%s)", url, url),
					},
					"logprobs":      nil,
					"finish_reason": "stop",
				},
			},
			"usage": gin.H{
				"prompt_tokens":     0,
				"completion_tokens": 0,
				"total_tokens":      0,
			},
		}

		if req.Stream {
			c.SSEvent("image", ret)
		} else {

			c.JSON(http.StatusOK, ret)
		}
	})

	r.Run()
}
