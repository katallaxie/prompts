package perplexity

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"iter"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/prompts"
)

// DefaultURL is the default endpoint for the Perplexity API.
const DefaultURL = "https://api.perplexity.ai/chat/completions"

// DefaultModel is the default model for the Perplexity API.
const DefaultModel = "sonar-pro"

// Defaults returns the default options for the Perplexity API.
func Defaults() []prompts.Opt[Event] {
	return []prompts.Opt[Event]{
		prompts.WithURL[Event](DefaultURL),
		prompts.WithTransformer(Transformer),
		prompts.WithDecoder(Decoder),
	}
}

// Event is the structure of the event received from the server.
type Event struct {
	// Type is the type of the event.
	Type string `json:"type"`
	// Data is the data of the event.
	Data []byte `json:"data"`
}

// Transformer is a function that transforms an event into a response.
var Transformer = func(e Event) (*prompts.ChatCompletionResponse, error) {
	resp := &prompts.ChatCompletionResponse{}

	if err := json.Unmarshal(e.Data, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

const maxBufferSize = 512 * 1 * 1000

// Decoder is an interface that defines the methods for decoding events from the response body.
var Decoder = func(body io.ReadCloser) iter.Seq[Event] {
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
	}
}
