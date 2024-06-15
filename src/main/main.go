package main

import (
	"api-autoMakeHtml/src/chat"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// コマンドライン引数から質問テキストを取得する
	if len(os.Args) < 2 {
		panic("too few arguments")
	}
	content := os.Args[1]

	err := godotenv.Load(".env")
	
	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 
	// 環境変数からAPIキーを取得する
	secret:= os.Getenv("OPEN_AI_SECRET")

	// リソース節約のためにタイムアウトを設定する
	timeout := 15 * time.Second

	// トークン節約のために応答の最大トークンを設定する
	maxTokens := 500

	// チャットに使用するモデルのID
	// modelID := "gpt-3.5-turbo"
	modelID := "gpt-4o"

	c := chat.NewChatCompletions(modelID, secret, maxTokens, timeout)
	res, err := c.AskOneQuestion(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("In %d / Out %d / Total %d tokens\n", res.Usage.PromptTokens, res.Usage.CompletionTokens, res.Usage.TotalTokens)
	for _, v := range res.Choices {
		fmt.Printf("[%s]: %s\n", v.Message.Role, v.Message.Content)
	}

}
