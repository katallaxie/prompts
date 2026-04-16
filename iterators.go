package prompts

import (
	"fmt"

	"github.com/katallaxie/pkg/logx"
)

// Log is a function that logs the event.
func Log(gen Generator) error {
	for msg, err := range gen {
		if err != nil {
			return err
		}

		logx.Infow("Received message", map[string]interface{}{
			"message": msg,
		})
	}

	return nil
}

// Print is a function that prints the event.
func Print(gen Generator) error {
	for msg, err := range gen {
		if err != nil {
			return err
		}

		fmt.Print(msg)
	}

	return nil
}
