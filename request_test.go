package prompts_test

import (
	"testing"

	"github.com/katallaxie/prompts"
	"github.com/stretchr/testify/require"
)

func TestNewChatCompletionRequest(t *testing.T) {
	msgs := []prompts.ChatCompletionMessage{
		{
			Role:    prompts.RoleUser,
			Content: "Hello, world!",
		},
	}

	req := prompts.NewChatCompletionRequest(msgs...)
	require.NotNil(t, req)
	require.Equal(t, msgs, req.Messages)
}

func TestNewStreamChatCompletionRequest(t *testing.T) {
	msgs := []prompts.ChatCompletionMessage{
		{
			Role:    prompts.RoleUser,
			Content: "Hello, world!",
		},
	}

	req := prompts.NewStreamChatCompletionRequest(msgs...)
	require.NotNil(t, req)
	require.Equal(t, msgs, req.Messages)
	require.True(t, req.Stream)
}
