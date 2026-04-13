package prompts

import (
	"fmt"
	"iter"

	"github.com/katallaxie/pkg/logx"
)

// Log is a function that logs the event.
func Log(res *ChatCompletionResponse) error {
	logx.Infow("Received response", map[string]interface{}{
		"model":   res.Model,
		"choices": res.Choices,
	})

	return nil
}

// Print is a function that prints the event.
func Print(stream iter.Seq2[*ChatCompletionResponse, error]) error {
	for msg, err := range stream {
		if err != nil {
			return err
		}

		fmt.Print(msg)
	}

	return nil
}
