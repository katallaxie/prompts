package perplexity

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/pkg/slices"
	"github.com/katallaxie/prompts"
)

// PerplexityChatCompletionRequest is the structure of the chat completion request for the Perplexity API.
type PerplexityChatCompletionRequest struct {
	// Model is the model for the chat completion request.
	Model string `json:"model"`
	// Messages is the list of messages for the chat completion request.
	Messages []struct {
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text,omitempty"`
		} `json:"content"`
	} `json:"messages"`
	// MaxTokens is the maximum number of tokens for the chat completion request.
	MaxTokens *int `json:"max_tokens,omitzero"`
	// Temperature is the sampling temperature
	Temperature *float32 `json:"temperature,omitzero"`
	// Stream is a flag to enable streaming
	Stream bool `json:"stream,omitempty"`
	// TopP is the nucleus sampling parameter
	TopP *float64 `json:"top_p,omitzero"`
	// TopK is the number of top tokens to sample from
	TopK *int `json:"top_k,omitzero"`
}

// UnmarshalPrompt unmarshals a ChatCompletionRequest into a PerplexityChatCompletionRequest.
func (p *PerplexityChatCompletionRequest) UnmarshalPrompt(in *prompts.ChatCompletionRequest) error {
	req := &PerplexityChatCompletionRequest{
		Model:       in.Model,
		MaxTokens:   in.MaxTokens,
		Temperature: in.Temperature,
		Stream:      in.Stream,
		TopP:        in.TopP,
		TopK:        in.TopK,
	}

	for _, msg := range in.Messages {
		m := struct {
			Role    string `json:"role"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text,omitempty"`
			} `json:"content"`
		}{
			Role: msg.Role.String(),
			Content: []struct {
				Type string `json:"type"`
				Text string `json:"text,omitempty"`
			}{},
		}

		for _, c := range msg.Content {
			if textContent, ok := c.GetText(); ok {
				m.Content = append(m.Content, struct {
					Type string `json:"type"`
					Text string `json:"text,omitempty"`
				}{
					Type: "text",
					Text: textContent.Text,
				})
			}
		}

		req.Messages = append(req.Messages, m)
	}

	p.Model = req.Model
	p.Messages = req.Messages
	p.MaxTokens = req.MaxTokens
	p.Temperature = req.Temperature
	p.Stream = req.Stream
	p.TopP = req.TopP
	p.TopK = req.TopK

	return nil
}

const maxBufferSize = 512 * 1 * 1000

// DefaultURL is the default endpoint for the Perplexity API.
const DefaultURL = "https://api.perplexity.ai/chat/completions"

// DefaultModel is the default model for the Perplexity API.
const DefaultModel = "sonar-pro"

// Defaults returns the default options for the Perplexity API.
func Defaults(opts ...prompts.Opt) []prompts.Opt {
	defaults := []prompts.Opt{
		prompts.WithURL(DefaultURL),
		prompts.WithClient(http.DefaultClient),
	}

	return slices.Append(defaults, opts...)
}

// Event is the structure of the event received from the server.
type Event struct {
	// Type is the type of the event.
	Type string `json:"type"`
	// Data is the data of the event.
	Data []byte `json:"data"`
}

var _ prompts.Prompter = (*Perplexity)(nil)

// Perplexity is a prompter that implements the Prompter interface for the Perplexity API.
type Perplexity struct{}

var _ prompts.StreamTransformer[Event] = (*Transformer)(nil)

// Transformer is a struct that implements the StreamTransformer interface for the Perplexity API.
type Transformer struct{}

// Transform transforms an event into a ChatCompletionResponse.
func (t *Transformer) Transform(iter iter.Seq[Event]) prompts.Stream {
	return func(yield func(*prompts.ChatCompletionResponse, error) bool) {
		for e := range iter {
			var res prompts.ChatCompletionResponse
			if err := json.Unmarshal(e.Data, &res); err != nil {
				if !yield(nil, err) {
					break
				}
				continue
			}

			if !yield(&res, nil) {
				break
			}
		}
	}
}

// NewTransformer creates a new Transformer.
func NewTransformer() *Transformer {
	return &Transformer{}
}

var _ prompts.StreamDecoder[Event] = (*Decoder)(nil)

// Decoder is a struct that implements the StreamDecoder interface for the Perplexity API.
type Decoder struct{}

// Decode decodes the response body into a stream of events.
func (d *Decoder) Decode(body io.ReadCloser) iter.Seq[Event] {
	scn := bufio.NewScanner(body)
	scn.Split(bufio.ScanLines)
	scn.Buffer(make([]byte, maxBufferSize), maxBufferSize)

	return func(yield func(Event) bool) {
		for scn.Scan() {
			event := cast.Zero[Event]()

			b := scn.Bytes()

			name, value, _ := bytes.Cut(b, []byte(":"))
			if len(value) > 0 && value[0] == ' ' {
				value = value[1:]
			}

			switch string(name) {
			case "":
				continue
			case "event":
				event.Type = string(value)
			case "data":
				event.Data = value
			}

			if !yield(event) {
				break
			}
		}

		body.Close()
	}
}

// NewDecoder creates a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// New creates a new Perplexity prompter with the given options.
func New() *Perplexity {
	return &Perplexity{}
}

// SendCompletionRequest sends a chat completion request to the Perplexity API and returns the response.
func (p *Perplexity) SendCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest) (*prompts.ChatCompletionResponse, error) {
	prompt := &PerplexityChatCompletionRequest{}
	err := prompt.UnmarshalPrompt(req)
	if err != nil {
		return nil, err
	}

	res := &prompts.ChatCompletionResponse{}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, req.Opts.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+req.Opts.ApiKey)
	r.Header.Set("Accept", "application/json")

	resp, err := req.Opts.Client.Do(r)
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
func (p *Perplexity) SendStreamCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest) (prompts.Stream, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, req.Opts.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+req.Opts.ApiKey)
	r.Header.Set("Accept", "text/event-stream")
	r.Header.Set("Connection", "keep-alive")

	resp, err := req.Opts.Client.Do(r) //nolint:bodyclose
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var promptErr prompts.PromptError
		err := json.NewDecoder(resp.Body).Decode(&promptErr)
		if err != nil {
			return nil, err
		}

		return nil, &promptErr
	}

	return NewTransformer().Transform(NewDecoder().Decode(resp.Body)), nil
}
