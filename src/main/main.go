package main

import (
	"api-autoMakeHtml/src/chat"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Message string `json:"message"`
}

func main() {
	engine := gin.Default()
	engine.POST("/submit", func(c *gin.Context) {
		// リクエストをバインドする
		var req RequestBody
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
	engine.Run(":8080")
}
