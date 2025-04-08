package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/prompts"
	"github.com/ollama/ollama/format"
)

// StreamResponse is the response from the Ollama API.
//
// Example:
//
//	{
//			"model": "llama3.2",
//			"created_at": "2023-08-04T08:52:19.385406455-07:00",
//			"message": {
//			"role": "assistant",
//	 	"content": "The",
//	 	"images": null
//		},
//		"done": false
//	}
type StreamResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
		Images  []struct {
			URL  string `json:"url"`
			Size int    `json:"size"`
		} `json:"images"`
	} `json:"message"`
	Done bool `json:"done"`
}

// DefaultURL is the default endpoint for the Ollama API.
const DefaultURL = "http://localhost:7869/api/chat"

// DefaultTimeout is the default timeout for the Ollama API.
const DefaultTimeout = 30 * time.Second

// DefaultModel is the default model for the Ollama API.
const DefaultModel = "smollm"

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

// Ollama is a chat model.
type Ollama struct {
	opts *Opts
}

// New returns a new Ollama.
func New(opts ...Opt) *Ollama {
	options := Defaults()

	for _, opt := range opts {
		opt(options)
	}

	p := new(Ollama)
	p.opts = options

	return p
}

// SendCompletionRequest sends a completion request to the Perplexity API.
func (o *Ollama) SendCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest) (*prompts.ChatCompletionResponse, error) {
	res := &prompts.ChatCompletionResponse{}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, o.opts.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	resp, err := o.opts.Client.Do(r)
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

const maxBufferSize = 512 * format.KiloByte

// SendStreamCompletionRequest sends a streamed completion request to the Ollama.
// nolint:gocyclo
func (o *Ollama) SendStreamCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest, res chan<- *prompts.ChatCompletionResponse) error {
	defer close(res)

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, o.opts.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "text/event-stream")
	r.Header.Set("Connection", "keep-alive")

	resp, err := o.opts.Client.Do(r)
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

	scanner := bufio.NewScanner(resp.Body)
	scanBuf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(scanBuf, maxBufferSize)

	for scanner.Scan() {
		stream := &StreamResponse{}

		bts := scanner.Bytes()
		if err := json.Unmarshal(bts, &stream); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		if len(bts) == 0 {
			continue
		}

		msg := &prompts.ChatCompletionResponse{
			Model: stream.Model,
			Choices: []prompts.ChatCompletionChoice{
				{
					Message: prompts.Index{
						Role:    prompts.Role(stream.Message.Role),
						Content: stream.Message.Content,
					},
				},
			},
		}

		res <- msg
	}

	return nil
}
