package callbacks

import (
	"github.com/katallaxie/pkg/logx"
	"github.com/katallaxie/prompts"
)

// Log is a function that logs the event.
func Log(res *prompts.ChatCompletionResponse) error {
	logx.Infow("Received response", map[string]interface{}{
		"model":   res.Model,
		"choices": res.Choices,
	})

	return nil
}
