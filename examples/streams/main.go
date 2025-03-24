package main

import (
	"context"

	"github.com/katallaxie/pkg/slices"
	"github.com/katallaxie/prompts"
	"github.com/katallaxie/prompts/ollama"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
)

const (
	example = `Write a concise summary of the following in less then 20 words:
	
	"Artificial intelligence (AI) is technology that enables computers and machines to simulate human learning, comprehension, problem solving, decision making, creativity and autonomy."
	
	CONCISE SUMMARY:`
)

func mapCompletionMessages(msg prompts.Completion) string {
	f := slices.First(msg.Choices...)
	return f.Message.GetContent()
}

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

	source := sources.NewChanSource(res)
	sink := sinks.NewStdout()

	source.Pipe(streams.NewPassThrough()).Pipe(streams.NewMap(mapCompletionMessages)).To(sink)
}
