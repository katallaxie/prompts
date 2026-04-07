package prompts

import "context"

// Prompt is a prompt that can be sent to a model.
type Prompt interface {
	// SendCompletionRequest sends a chat completion request.
	SendCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	// SendStreamCompletionRequest sends a chat completion request and streams the response.
	SendStreamCompletionRequest(ctx context.Context, req *ChatCompletionRequest, cb ...func(res *ChatCompletionResponse) error) error
}
