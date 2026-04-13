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
	prompt := prompts.New(
		prompts.WithDecoder(perplexity.Decoder),
		prompts.WithTransformer(perplexity.Transformer),
		prompts.WithURL[perplexity.Event](perplexity.DefaultURL),
		prompts.WithApiKey[perplexity.Event](os.Getenv("PPLX_API_KEY")),
	)

	msgs := []prompts.ChatCompletionMessage{
		{
			Role:    prompts.RoleSystem,
			Content: "You are a helpful assistant. You start every answer with 'Sure my lord!'",
		},
		{
			Role:    prompts.RoleUser,
			Content: "What is the definition of Pi?",
		},
	}

	req := prompts.NewStreamChatCompletionRequest(msgs...)
	req.SetModel(perplexity.DefaultModel)

	err := prompt.SendStreamCompletionRequest(context.Background(), req, prompts.Print)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
