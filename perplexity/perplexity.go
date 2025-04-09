package perplexity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/prompts"
)

// DefaultURL is the default endpoint for the Perplexity API.
const DefaultURL = "https://api.perplexity.ai/chat/completions"

// DefaultTimeout is the default timeout for the Perplexity API.
const DefaultTimeout = 30 * time.Second

// DefaultModel is the default model for the Perplexity API.
const DefaultModel = "sonar-pro"

// NewChatCompletionRequest creates a new chat completion request
func NewChatCompletionRequest() *prompts.ChatCompletionRequest {
	return &prompts.ChatCompletionRequest{
		Model:    DefaultModel,
		Messages: []prompts.ChatCompletionMessage{},
		Stream:   cast.Ptr(false),
	}
}

// NewStreamCompletionRequest creates a new chat stream completion request
func NewStreamCompletionRequest() *prompts.ChatCompletionRequest {
	return &prompts.ChatCompletionRequest{
		Model:    DefaultModel,
		Messages: []prompts.ChatCompletionMessage{},
		Stream:   cast.Ptr(true),
	}
}

// Opts ...
type Opts struct {
	// BaseURL is the base URL.
	BaseURL string `json:"base_url"`
	// ApiKey is the API key.
	ApiKey string `json:"api_key"`
	// Timeout is the timeout.
	Timeout time.Duration `json:"timeout"`
	// Client is the HTTP client.
	Client *http.Client `json:"-"`
}

// Opt ...
type Opt func(*Opts)

// WithURL configures the base URL.
func WithURL(url string) Opt {
	return func(o *Opts) {
		o.BaseURL = url
	}
}

// WithApiKey configures the API key.
func WithApiKey(apiKey string) Opt {
	return func(o *Opts) {
		o.ApiKey = apiKey
	}
}

// WithClient configures the HTTP client.
func WithClient(client *http.Client) Opt {
	return func(o *Opts) {
		o.Client = client
	}
}

// WithTimeout configures the timeout.
func WithTimeout(timeout time.Duration) Opt {
	return func(o *Opts) {
		o.Client.Timeout = timeout
	}
}

// WithBaseURL configures the base URL.
func WithBaseURL(url string) Opt {
	return func(o *Opts) {
		o.BaseURL = url
	}
}

// Defaults is the default options.
func Defaults() *Opts {
	return &Opts{
		BaseURL: DefaultURL,
		Timeout: DefaultTimeout,
		Client:  &http.Client{},
	}
}

// Perplexity is a chat model.
type Perplexity struct {
	opts *Opts
}

var _ prompts.Chat = (*Perplexity)(nil)

// New returns a new Perplexity.
func New(opts ...Opt) *Perplexity {
	options := Defaults()

	for _, opt := range opts {
		opt(options)
	}

	p := new(Perplexity)
	p.opts = options

	return p
}

// SendCompletionRequest sends a completion request to the Perplexity API.
func (p *Perplexity) SendCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest) (*prompts.ChatCompletionResponse, error) {
	res := &prompts.ChatCompletionResponse{}

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

// SendStreamCompletionRequest sends a streamed completion request to the Perplexity API.
// nolint:gocyclo
func (p *Perplexity) SendStreamCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest, cb ...func(res *prompts.ChatCompletionResponse) error) error {
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return prompts.FromBody(body)
	}

	stream := prompts.NewStream(NewDecoder(resp), Transformer)
	defer stream.Close()

	return prompts.Events(stream.All(), cb...)
}
