package main

import (
	"context"
	"fmt"
	chatGpt "github.com/sashabaranov/go-gpt3"
)

func main() {
	title := "帮我写一篇大学生毕业论文"
	toKen := "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"
	c := chatGpt.NewClient(toKen)
	ctx := context.Background()
	req := chatGpt.CompletionRequest{
		Model:            chatGpt.GPT3TextDavinci003,
		MaxTokens:        4000,
		Temperature:      0.5,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Prompt:           title,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("↓↓")
	fmt.Println("回复的内容：", resp.Choices[0].Text)
}
