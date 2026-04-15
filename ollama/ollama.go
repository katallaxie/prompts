package ollama

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

const maxBufferSize = 512 * 1 * 1000

// DefaultURL is the default endpoint for the Ollama API.
const DefaultURL = "http://localhost:7869/api/chat"

// DefaultModel is the default model for the Ollama API.
const DefaultModel = "smollm"

// Defaults returns the default options for the Ollama API.
func Defaults(opts ...prompts.Opt) []prompts.Opt {
	defaults := []prompts.Opt{
		prompts.WithURL(DefaultURL),
		prompts.WithClient(http.DefaultClient),
	}

	return slices.Append(defaults, opts...)
}

// Event is the structure of the event stream response from the Ollama API.
type Event struct {
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

var _ prompts.StreamTransformer[Event] = (*Transformer)(nil)

// Transformer is a struct that implements the StreamTransformer interface for the Ollama API.
type Transformer struct{}

// Transform transforms an event into a ChatCompletionResponse.
func (t *Transformer) Transform(iter iter.Seq[Event]) prompts.Stream {
	return func(yield func(*prompts.ChatCompletionResponse, error) bool) {
		for e := range iter {
			var res prompts.ChatCompletionResponse
			res.Model = e.Model
			res.Choices = []prompts.ChatCompletionChoice{
				{
					Message: prompts.ChatCompletionChoiceIndex{
						Role:    prompts.Role(e.Message.Role),
						Content: e.Message.Content,
					},
				},
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

// Decoder is a struct that implements the StreamDecoder interface for the Ollama API.
type Decoder struct{}

// Decode decodes the event stream response from the Ollama API into a sequence of Events.
func (d *Decoder) Decode(body io.ReadCloser) iter.Seq[Event] {
	scn := bufio.NewScanner(body)
	scn.Split(bufio.ScanLines)
	scn.Buffer(make([]byte, maxBufferSize), maxBufferSize)

	return func(yield func(Event) bool) {
		for scn.Scan() {
			event := cast.Zero[Event]()

			b := scn.Bytes()
			if len(b) == 0 {
				continue
			}

			if err := json.Unmarshal(b, &event); err != nil {
				break
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

// Ollama is a prompter that implements the Prompter interface for the Ollama API.
type Ollama struct{}

var _ prompts.Prompter = (*Ollama)(nil)

// New creates a new Ollama prompter with the given options.
func New() *Ollama {
	return &Ollama{}
}

// SendCompletionRequest sends a chat completion request to the Ollama API and returns the response.
func (p *Ollama) SendCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest) (*prompts.ChatCompletionResponse, error) {
	res := &prompts.ChatCompletionResponse{}

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
func (p *Ollama) SendStreamCompletionRequest(ctx context.Context, req *prompts.ChatCompletionRequest) (prompts.Stream, error) {
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
