package main

import (
	"context"
	"fmt"
	"os"

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
			Content: "You are a helpful assistant. You start every answers with 'Sure!'",
		},
		{
			Role:    "user",
			Content: "What is the definition of Pi?",
		},
	}

	req := prompts.NewStreamChatCompletionRequest()
	req.SetModel(perplexity.DefaultModel)
	req.AddMessages(msg...)

	err := client.SendStreamCompletionRequest(context.Background(), req, prompts.Print)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
