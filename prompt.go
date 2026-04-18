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

// DefaultClient is the default HTTP client for the Perplexity API.
var DefaultClient = &http.Client{
	Timeout: DefaultTimeout,
}

// Generator is a type that represents a generator of chat completion responses.
// It is an iterator that yields chat completion responses or errors.
type Generator iter.Seq2[*Response, error]

// Decoder is a function type that implements the Decoder interface.
type Decoder[E any] interface {
	Decode(io.ReadCloser) iter.Seq[E]
}

// Transformer is a function type that implements the Transformer interface.
type Transformer[E any] interface {
	Transform(iter.Seq[E]) Generator
}

// Responder is the interface for sending a chat completion request and receiving a response.
type Responder interface {
	// CreateResponse sends a chat completion request and returns the response.
	CreateResponse(ctx context.Context, req *ResponseRequest) (*Response, error)
}

// Client is the main struct for interacting with the Perplexity API.
type Client struct {
	Opts     *Opts
	Response Responder
}

// New creates a Promts client with the given options.
func New(opts ...Opt) *Client {
	options := new(Opts)
	options.Client = DefaultClient

	client := new(Client)
	client.Opts = options

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// ResponderFactory is a function type for creating a responder.
type ResponderFactory func(c *Client) Responder

// Opt is a function type for configuring the Client.
type Opt func(*Client)

// Opts is a function that applies multiple options to the Client.
type Opts struct {
	// BaseURL is the base URL.
	BaseURL string `json:"base_url"`
	// ApiKey is the API key.
	ApiKey string `json:"api_key"`
	// Headers are the headers to include in the request.
	Headers map[string]string `json:"headers"`
	// Client is the HTTP client.
	Client *http.Client `json:"-"`
}

// WithURL configures the base URL.
func WithURL(url string) Opt {
	return func(c *Client) {
		c.Opts.BaseURL = url
	}
}

// WithApiKey configures the API key.
func WithApiKey(apiKey string) Opt {
	return func(c *Client) {
		c.Opts.ApiKey = apiKey
	}
}

// WithClient configures the HTTP client.
func WithClient(client *http.Client) Opt {
	return func(c *Client) {
		c.Opts.Client = client
	}
}

// WithBaseURL configures the base URL.
func WithBaseURL(url string) Opt {
	return func(c *Client) {
		c.Opts.BaseURL = url
	}
}

// WithResponder configures the response for the client.
func WithResponder(factory ResponderFactory) Opt {
	return func(c *Client) {
		client := factory(c)
		c.Response = client
	}
}
