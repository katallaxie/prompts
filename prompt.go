package prompts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultTimeout is the default timeout for the Perplexity API.
const DefaultTimeout = 30 * time.Second

// Prompter is the interface for sending chat completion requests to a language model.
// It abstracts away the details of how the requests are sent and allows for different implementations (e.g., using different APIs or libraries).
type Prompter interface {
	// SendCompletionRequest sends a chat completion request.
	SendCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	// SendStreamCompletionRequest sends a chat completion request and streams the response.
	SendStreamCompletionRequest(ctx context.Context, req *ChatCompletionRequest, iter StreamIterator) error
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
	// Decoder is the decoder for the response.
	Decoder StreamDecoder[E] `json:"-"`
	// Transformer is the transformer for the response.
	Transformer StreamTransformer[E] `json:"-"`
}

// Prompt is the main struct for sending chat completion requests.
type Prompt[E any] struct {
	opts *Opts[E]
}

var _ Prompter = (*Prompt[any])(nil)

// New returns a new Prompt.
func New[E any](opts ...Opt[E]) *Prompt[E] {
	options := new(Opts[E])
	options.Client = &http.Client{Timeout: DefaultTimeout}

	for _, opt := range opts {
		opt(options)
	}

	p := new(Prompt[E])
	p.opts = options

	return p
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

// WithTransformer configures the transformer for the response.
func WithTransformer[E any](transformer StreamTransformer[E]) Opt[E] {
	return func(o *Opts[E]) {
		o.Transformer = transformer
	}
}

// WithDecoder configures the decoder for the response.
func WithDecoder[E any](decoder StreamDecoder[E]) Opt[E] {
	return func(o *Opts[E]) {
		o.Decoder = decoder
	}
}

// SendCompletionRequest sends a chat completion request.
func (p *Prompt[E]) SendCompletionRequest(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	res := &ChatCompletionResponse{}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, p.opts.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+p.opts.ApiKey)
	r.Header.Set("Accept", "application/json")

	resp, err := p.opts.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SendStreamCompletionRequest sends a chat completion request and streams the response.
func (p *Prompt[E]) SendStreamCompletionRequest(ctx context.Context, req *ChatCompletionRequest, iter StreamIterator) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, p.opts.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+p.opts.ApiKey)
	r.Header.Set("Accept", "text/event-stream")
	r.Header.Set("Connection", "keep-alive")

	resp, err := p.opts.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var promptErr PromptError
		err := json.NewDecoder(resp.Body).Decode(&promptErr)
		if err != nil {
			return err
		}

		return &promptErr
	}

	stream := NewStream(resp.Body, p.opts.Decoder, p.opts.Transformer)

	err = iter(stream.All())
	if err != nil {
		return err
	}

	return nil
}
