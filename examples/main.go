package main

import (
	"context"
	"encoding/json"
	"fmt"
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
						Text: "What is my horoscope? I am an Aquarius.",
					},
				},
			},
		},
	}

	req := prompts.NewChatCompletionRequest(perplexity.Defaults(prompts.WithApiKey(os.Getenv("PPLX_API_KEY")), prompts.WithInput(msgs...))...)
	req.Model = perplexity.DefaultModel
	req.Tools = []prompts.ChatCompletionTool{
		{
			Tool: prompts.ChatCompletionFuntionTool{
				Function: prompts.ChatCompletionFunctionDefintion{
					Name:        "get_horoscope",
					Description: "Get the horoscope for a given zodiac sign.",
					Parameters: prompts.ChatCompletionFunctionParameters{
						Properties: prompts.ChatCompletionFunctionProperties{
							"sign": json.RawMessage(`{"type": "string", "enum": ["aries", "taurus", "gemini", "cancer", "leo", "virgo", "libra", "scorpio", "sagittarius", "capricorn", "aquarius", "pisces"]}`),
						},
						Required: []string{"sign"},
					},
				},
			},
		},
	}

	res, err := client.SendCompletionRequest(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Output)
}
