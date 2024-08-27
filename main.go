package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/k8scat/flux-api/api"
	"github.com/k8scat/flux-api/translate"
)

func translatePrompt(prompt string) string {
	enable := strings.ToLower(os.Getenv("TRANSLATE_ENABLE"))
	if enable == "true" {
		translator := translate.NewOpenAITranslator(
			os.Getenv("TRANSLATE_BASE_URL"),
			os.Getenv("TRANSLATE_API_KEY"),
			os.Getenv("TRANSLATE_MODEL"),
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

func main() {
	r := gin.Default()

	r.POST("/v1/images/generations", func(c *gin.Context) {
		var req api.GenerateImageRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Model != "FLUX.1-schnell" && req.Model != "flux-v1" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid model"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		auth := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if auth == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auth"})
			return
		}

		req.Prompt = translatePrompt(req.Prompt)

		var client api.OpenAI
		if req.Model == "FLUX.1-schnell" {
			client = api.NewSiliconFlow(auth)
		} else if req.Model == "flux-v1" {
			client = api.NewGetimgAI(auth)
		}

		resp, err := client.GenerateImage(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.Run()
}
