package prompts

import (
	"context"
	"io"
	"iter"
	"time"
)

// DefaultTimeout is the default timeout for the Perplexity API.
const DefaultTimeout = 30 * time.Second

// Generator is a type that represents a generator of chat completion responses.
// It is an iterator that yields chat completion responses or errors.
type Generator iter.Seq2[*ChatCompletionResponse, error]

// Decoder is a function type that implements the Decoder interface.
type Decoder[E any] interface {
	Decode(io.ReadCloser) iter.Seq[E]
}

// Transformer is a function type that implements the Transformer interface.
type Transformer[E any] interface {
	Transform(iter.Seq[E]) Generator
}

// Prompter is the interface for sending chat completion requests to a language model.
// It abstracts away the details of how the requests are sent and allows for different implementations (e.g., using different APIs or libraries).
type Prompter interface {
	// SendCompletionRequest sends a chat completion request.
	SendCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	// SendStreamCompletionRequest sends a chat completion request and streams the response.
	SendStreamCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (Generator, error)
}
