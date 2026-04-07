package callbacks

import (
	"fmt"

	"github.com/katallaxie/prompts"
)

// Print is a function that prints the event.
func Print(res *prompts.ChatCompletionResponse) error {
	fmt.Print(res)

	return nil
}
