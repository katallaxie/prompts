package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/katallaxie/prompts"
	"github.com/katallaxie/prompts/openai"
	"github.com/katallaxie/prompts/perplexity"
)

// This example demonstrates how to create a completion request with a message
// It then sends the request to the API and prints the last completion content.
func main() {
	client := prompts.NewClient().APIKey(os.Getenv("PPLX_API_KEY"))
	prompt := perplexity.New(client)

	msgs := []perplexity.ResponseInput{
		{
			Role: perplexity.RoleSystem,
			Content: []perplexity.ResponseMessageContent{
				{
					Content: perplexity.ResponseMessageContentText{
						Text: "You are a helpful assistant. You answer questions to the best of your ability.",
					},
				},
			},
		},
		{
			Role: perplexity.RoleUser,
			Content: []perplexity.ResponseMessageContent{
				{
					Content: perplexity.ResponseMessageContentText{
						Text: "What is my horoscope? I am an Aquarius.",
					},
				},
			},
		},
	}

	req := openai.NewResponseRequest(openai.WithInput(msgs...))
	req.Model = perplexity.DefaultModel
	req.Tools = []perplexity.ResponseTool{
		{
			Tool: perplexity.ResponseFunctionTool{
				Function: perplexity.ResponseFunctionDefinition{
					Name:        "get_horoscope",
					Description: "Get the horoscope for a given zodiac sign.",
					Parameters: perplexity.ResponseFunctionParameters{
						Properties: perplexity.ResponseFunctionProperties{
							"sign": json.RawMessage(`{"type": "string", "enum": ["aries", "taurus", "gemini", "cancer", "leo", "virgo", "libra", "scorpio", "sagittarius", "capricorn", "aquarius", "pisces"]}`),
						},
						Required: []string{"sign"},
					},
				},
			},
		},
	}

	res, err := prompt.Respond(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Output)
}
