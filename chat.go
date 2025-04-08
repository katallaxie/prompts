package prompts

import "context"

// Chat is an interface for chat completions.
type Chat interface {
	// SendCompletionRequest sends a chat completion request.
	SendCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	// SendStreamCompletionRequest sends a chat completion request and streams the response.
	SendStreamCompletionRequest(ctx context.Context, req *ChatCompletionRequest, cb ...func(res *ChatCompletionResponse) error) error
}
