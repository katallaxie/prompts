package main

import (
	"context"
	"fmt"
	"os"

	"github.com/katallaxie/prompts"
	"github.com/katallaxie/prompts/ollama"
)

// This example demonstrates how to create a completion request with a message
// It then sends the request to the API and prints the last completion content.
func main() {
	client := ollama.New()
	msg := []prompts.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant. You start every answers with 'Sure!'",
		},
		{
			Role:    "user",
			Content: "What is the definition of Pi?",
		},
	}

	req := ollama.NewStreamCompletionRequest()
	req.AddMessages(msg...)

	stream := make(chan *prompts.ChatCompletionResponse, 1)

	go func() {
		for msg := range stream {
			fmt.Print(msg)
		}
	}()

	err := client.SendStreamCompletionRequest(context.Background(), req, stream)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
