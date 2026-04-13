package ollama

import (
	"bufio"
	"encoding/json"
	"io"
	"iter"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/prompts"
)

// DefaultURL is the default endpoint for the Ollama API.
const DefaultURL = "http://localhost:7869/api/chat"

// DefaultModel is the default model for the Ollama API.
const DefaultModel = "smollm"

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

// Transformer is a function that transforms an Event into a ChatCompletionResponse.
var Transformer = func(e Event) (*prompts.ChatCompletionResponse, error) {
	return &prompts.ChatCompletionResponse{
		Model: e.Model,
		Choices: []prompts.ChatCompletionChoice{
			{
				Message: prompts.Index{
					Role:    prompts.Role(e.Message.Role),
					Content: e.Message.Content,
				},
			},
		},
	}, nil
}

const maxBufferSize = 512 * 1 * 1000

// Decoder is a function that decodes the event stream response from the Ollama API into a sequence of Events.
var Decoder = func(body io.ReadCloser) iter.Seq[Event] {
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
	}
}
