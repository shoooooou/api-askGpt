package icon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const openaiURL = "https://api.openai.com/v1/images/generations"

func GenerateImage(prompt, size string) (string, error) {
	client := resty.New()
	err := godotenv.Load(".env")
	// 環境変数が読み込めない場合はエラーを返す
	if err != nil || os.Getenv("OPEN_AI_SECRET") == "" {
		panic("envFile not found or OPEN_AI_SECRET is empty")
	}
	response, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("OPEN_AI_SECRET")).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":           "dall-e-3",
			"prompt":          prompt,
			"n":               1,
			"size":            size,
			"response_format": "url",
		}).
		Post(openaiURL)

	if err != nil {
		return "", err
	}

	if response.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("failed to generate image: %s", response.Status())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return "", err
	}

	imageURL := result["data"].([]interface{})[0].(map[string]interface{})["url"].(string)

	return imageURL, nil
}
