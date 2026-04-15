package prompts

import (
	"context"
	"io"
	"iter"
	"time"
)

// DefaultTimeout is the default timeout for the Perplexity API.
const DefaultTimeout = 30 * time.Second

// Stream is the interface for a stream of events.
type Stream iter.Seq2[*ChatCompletionResponse, error]

// StreamDecoder is a function type that implements the StreamDecoder interface.
type StreamDecoder[E any] interface {
	Decode(io.ReadCloser) iter.Seq[E]
}

// StreamTransformer is a function type that implements the StreamTransformer interface.
type StreamTransformer[E any] interface {
	Transform(iter.Seq[E]) Stream
}

// Prompter is the interface for sending chat completion requests to a language model.
// It abstracts away the details of how the requests are sent and allows for different implementations (e.g., using different APIs or libraries).
type Prompter interface {
	// SendCompletionRequest sends a chat completion request.
	SendCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	// SendStreamCompletionRequest sends a chat completion request and streams the response.
	SendStreamCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (Stream, error)
}
