package prompts

import (
	"context"
	"io"
	"iter"
	"net/http"
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

// Opt is a function that configures the options for the Prompter.
type Opt[E any] func(*Opts[E])

// Opts is the options for configuring the Prompter.
type Opts[E any] struct {
	// BaseURL is the base URL.
	BaseURL string `json:"base_url"`
	// ApiKey is the API key.
	ApiKey string `json:"api_key"`
	// Timeout is the timeout.
	Timeout time.Duration `json:"timeout"`
	// Client is the HTTP client.
	Client *http.Client `json:"-"`
}

// WithURL configures the base URL.
func WithURL[E any](url string) Opt[E] {
	return func(o *Opts[E]) {
		o.BaseURL = url
	}
}

// WithApiKey configures the API key.
func WithApiKey[E any](apiKey string) Opt[E] {
	return func(o *Opts[E]) {
		o.ApiKey = apiKey
	}
}

// WithClient configures the HTTP client.
func WithClient[E any](client *http.Client) Opt[E] {
	return func(o *Opts[E]) {
		o.Client = client
	}
}

// WithTimeout configures the timeout.
func WithTimeout[E any](timeout time.Duration) Opt[E] {
	return func(o *Opts[E]) {
		o.Client.Timeout = timeout
	}
}

// WithBaseURL configures the base URL.
func WithBaseURL[E any](url string) Opt[E] {
	return func(o *Opts[E]) {
		o.BaseURL = url
	}
}
