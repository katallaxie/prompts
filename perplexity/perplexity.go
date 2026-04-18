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
	"github.com/katallaxie/prompts"
)

const maxBufferSize = 512 * 1 * 1000

// DefaultURL is the default endpoint for the Perplexity API.
const DefaultURL = "https://api.perplexity.ai/v1/responses"

// DefaultModel is the default model for the Perplexity API.
const DefaultModel = "anthropic/claude-opus-4-6"

// Event is the structure of the event received from the server.
type Event struct {
	// Type is the type of the event.
	Type string `json:"type"`
	// Data is the data of the event.
	Data []byte `json:"data"`
}

var _ prompts.Responder = (*Perplexity)(nil)

// Perplexity is a prompter that implements the Responder interface for the Perplexity API.
type Perplexity struct {
	c *prompts.Client
}

var _ prompts.Transformer[Event] = (*Transformer)(nil)

// Transformer is a struct that implements the Transformer interface for the Perplexity API.
type Transformer struct{}

// Transform transforms an event into a ChatCompletionResponse.
func (t *Transformer) Transform(iter iter.Seq[Event]) prompts.Generator {
	return func(yield func(*prompts.Response, error) bool) {
		for e := range iter {
			var res prompts.Response
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

var _ prompts.Decoder[Event] = (*Decoder)(nil)

// Decoder is a struct that implements the Decoder interface for the Perplexity API.
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
func NewResponder() prompts.ResponderFactory {
	return func(c *prompts.Client) prompts.Responder {
		return &Perplexity{c: c}
	}
}

// CreateResponse sends a chat completion request and returns the response.
func (p *Perplexity) CreateResponse(ctx context.Context, req *prompts.ResponseRequest) (*prompts.Response, error) {
	res := &prompts.Response{}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, p.c.Opts.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+p.c.Opts.ApiKey)
	r.Header.Set("Accept", "application/json")

	resp, err := p.c.Opts.Client.Do(r)
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

	fmt.Println(string(body))

	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
