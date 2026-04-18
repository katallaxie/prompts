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
	prompt := prompts.New(
		prompts.WithApiKey(os.Getenv("PPLX_API_KEY")),
		prompts.WithBaseURL(perplexity.DefaultURL),
		prompts.WithResponder(perplexity.NewResponder()),
	)

	msgs := []prompts.ResponseInput{
		{
			Role: prompts.RoleSystem,
			Content: []prompts.ResponseMessageContent{
				{
					Content: prompts.ResponseMessageContentText{
						Text: "You are a helpful assistant. You answer questions to the best of your ability.",
					},
				},
			},
		},
		{
			Role: prompts.RoleUser,
			Content: []prompts.ResponseMessageContent{
				{
					Content: prompts.ResponseMessageContentText{
						Text: "What is my horoscope? I am an Aquarius.",
					},
				},
			},
		},
	}

	req := prompts.NewResponseRequest(prompts.WithInput(msgs...))
	req.Model = perplexity.DefaultModel
	req.Tools = []prompts.ResponseTool{
		{
			Tool: prompts.ResponseFunctionTool{
				Function: prompts.ResponseFunctionDefinition{
					Name:        "get_horoscope",
					Description: "Get the horoscope for a given zodiac sign.",
					Parameters: prompts.ResponseFunctionParameters{
						Properties: prompts.ResponseFunctionProperties{
							"sign": json.RawMessage(`{"type": "string", "enum": ["aries", "taurus", "gemini", "cancer", "leo", "virgo", "libra", "scorpio", "sagittarius", "capricorn", "aquarius", "pisces"]}`),
						},
						Required: []string{"sign"},
					},
				},
			},
		},
	}

	res, err := prompt.Response.CreateResponse(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Output)
}
