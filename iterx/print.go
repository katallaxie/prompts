package iterx

import (
	"fmt"
	"iter"

	"github.com/katallaxie/prompts"
)

// Print is a function that prints the event.
func Print(stream iter.Seq2[*prompts.ChatCompletionResponse, error]) error {
	for msg, err := range stream {
		if err != nil {
			return err
		}

		fmt.Println(msg)
	}

	return nil
}
