package main

import (
	"api-autoMakeHtml/src/chat"
	"api-autoMakeHtml/src/icon"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TextRequest struct {
	Message string `json:"message"`
}
type ImageRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Size   string `json:"size" binding:"required"`
}

func main() {
	engine := gin.Default()
	engine.POST("/askgpt/text/", func(c *gin.Context) {
		// リクエストをバインドする
		var req TextRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		// メッセージが空の場合はエラーを返す
		if req.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "message is empty",
			})
		}

		timeout := 15 * time.Second
		maxTokens := 500
		modelID := "gpt-4o"

		chatCompletion := chat.NewChatCompletions(modelID, maxTokens, timeout)
		res, err := chatCompletion.AskOneQuestion(req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		fmt.Printf(res.Choices[0].Message.Content)

		c.JSON(http.StatusOK, gin.H{
			"message": res.Choices[0].Message.Content,
		})
	})

	engine.POST("/askgpt/icon/", func(c *gin.Context) {
		var req ImageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		imageURL, err := icon.GenerateImage(req.Prompt, req.Size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"image_url": imageURL})
	})

	engine.Run(":8080")
}
