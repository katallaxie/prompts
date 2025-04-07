package main

import (
	"context"
	"fmt"
	"os"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/prompts"
	"github.com/katallaxie/prompts/perplexity"
)

// This example demonstrates how to create a completion request with a message
// It then sends the request to the API and prints the last completion content.
func main() {
	client := perplexity.New(perplexity.WithApiKey(os.Getenv("PPLX_API_KEY")))
	msg := []prompts.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant. You use emojies to add some fun to your responses.",
		},
		{
			Role:    "user",
			Content: "What is the definition of Pi?",
		},
	}

	req := &prompts.ChatCompletionRequest{
		Model:     "sonar-pro",
		Messages:  msg,
		MaxTokens: 100,
		Stream:    cast.Ptr(true),
	}

	stream := make(chan *prompts.ChatCompletionResponse, 1)

	go func() {
		for msg := range stream {
			for _, choice := range msg.Choices {
				fmt.Print(choice)
			}
		}
	}()

	err := client.SendStreamCompletionRequest(context.Background(), req, stream)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
