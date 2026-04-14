package prompts_test

import (
	"bytes"
	"testing"

	"github.com/katallaxie/prompts"
	"github.com/stretchr/testify/require"
)

func TestNewChatCompletionResponse(t *testing.T) {
	choices := []prompts.ChatCompletionChoice{
		{
			Message: prompts.ChatCompletionChoiceIndex{
				Content: "Hello, how can I help you?",
			},
		},
	}

	res := prompts.NewChatCompletionResponse(choices...)
	require.NotNil(t, res)
	require.Equal(t, choices, res.Choices)

	var b bytes.Buffer
	_, err := b.ReadFrom(res)
	require.NoError(t, err)
	require.Equal(t, "Hello, how can I help you?", b.String())
}
