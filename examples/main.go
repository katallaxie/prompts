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
	client := perplexity.New()

	msgs := []prompts.ChatCompletionMessage{
		{
			Role: prompts.RoleSystem,
			Content: []prompts.ChatCompletionMessageContent{
				{
					Content: prompts.ChatCompletionMessageContentText{
						Text: "You are a helpful assistant. You answer questions to the best of your ability.",
					},
				},
			},
		},
		{
			Role: prompts.RoleUser,
			Content: []prompts.ChatCompletionMessageContent{
				{
					Content: prompts.ChatCompletionMessageContentText{
						Text: "What is the definition of Pi?",
					},
				},
			},
		},
	}

	req := prompts.NewStreamChatCompletionRequest(perplexity.Defaults(prompts.WithApiKey(os.Getenv("PPLX_API_KEY")), prompts.WithMessages(msgs...))...)
	req.Model = perplexity.DefaultModel

	res, err := client.SendStreamCompletionRequest(context.Background(), req)
	if err != nil {
		panic(err)
	}

	prompts.Print(res)
}
