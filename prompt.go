package prompts

import "context"

// Prompt is a prompt that can be sent to a model.
type Prompt interface {
	// SendCompletionRequest sends a chat completion request.
	SendCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	// SendStreamCompletionRequest sends a chat completion request and streams the response.
	SendStreamCompletionRequest(ctx context.Context, req *ChatCompletionRequest, cb ...func(res *ChatCompletionResponse) error) error
}

// PromptChannel is a prompt that can be sent to a model and receives responses through a channel.
type PromptChannel <-chan *ChatCompletionResponse

// Abortable is a prompt that can be aborted.
type Abortable interface {
	// Abort aborts the prompt.
	Abort() error
}

// Subscriable is a prompt that can be subscribed to.
type Subscriable interface {
	// Subscribe subscribes to the prompt and returns a channel that will receive the responses.
	Subscribe() (PromptChannel, error)
}
