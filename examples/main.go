package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/katallaxie/pkg/conv"
	"github.com/katallaxie/prompts"
	"github.com/katallaxie/prompts/ollama"
)

const (
	example = `Write a concise summary of the following in less then 20 words:
	
	"Artificial intelligence (AI) is technology that enables computers and machines to simulate human learning, comprehension, problem solving, decision making, creativity and autonomy."
	
	CONCISE SUMMARY:`
)

func main() {
	api, err := ollama.New(ollama.WithBaseURL("http://localhost:7869"), ollama.WithModel("smollm"))
	if err != nil {
		panic(err)
	}

	prompt := prompts.Prompt{
		Model: prompts.Model("smollm"),
		Messages: []prompts.Message{
			&prompts.UserMessage{
				Content: example,
			},
		},
	}

	res, err := api.Complete(context.Background(), &prompt)
	if err != nil {
		panic(err)
	}

	r := prompts.NewCompletionReader(res.Choices...)

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r)
	fmt.Print(conv.String(buf))
}
