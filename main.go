package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/k8scat/flux-api/api"
)

func main() {
	r := gin.Default()

	r.POST("/v1/images/generations", func(c *gin.Context) {
		var req api.GenerateImageRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authHeader := c.GetHeader("Authorization")
		auth := strings.TrimPrefix(authHeader, "Bearer ")

		var client api.OpenAI
		if req.Model == "FLUX.1-schnell" {
			client = api.NewSiliconFlow(auth)
		} else if req.Model == "flux-v1" {
			client = api.NewGetimgAI(auth)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
			return
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
