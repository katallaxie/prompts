package main

import (
	"context"
	"os"

	"github.com/katallaxie/prompts"
	"github.com/katallaxie/prompts/perplexity"
)

// This example demonstrates how to create a completion request with a message
// It then sends the request to the API and prints the last completion content.
func main() {
	client := perplexity.New(
		perplexity.Defaults(prompts.WithApiKey[perplexity.Event](os.Getenv("PPLX_API_KEY")))...)

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
	req.Model = perplexity.DefaultModel

	stream, err := client.SendStreamCompletionRequest(context.Background(), req)
	if err != nil {
		panic(err)
	}

	prompts.Print(stream)
}
