package ollama

import (
	"bufio"
	"encoding/json"
	"io"
	"iter"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/pkg/utilx"
	"github.com/katallaxie/prompts"
)

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

// NewDecoder creates a new Decoder based on the content type of the response.
func NewDecoder(body io.ReadCloser) prompts.StreamDecoder[Event] {
	if utilx.Empty(body) {
		return nil
	}

	scanner := bufio.NewScanner(body)
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)

	return &eventStreamDecoder[Event]{rc: body, scn: scanner}
}

type eventStreamDecoder[E any] struct {
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

// All returns an iterator over all events to be decoded.
func (s *eventStreamDecoder[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for s.scn.Scan() {
			event := cast.Zero[E]()

			b := s.scn.Bytes()
			if len(b) == 0 {
				continue
			}

			if err := json.Unmarshal(b, &event); err != nil {
				s.err = err
				break
			}

			if !yield(event) {
				break
			}
		}

		if s.scn.Err() != nil {
			s.err = s.scn.Err()
		}
	}
}

// Close closes the decoder.
func (s *eventStreamDecoder[E]) Close() error {
	return s.rc.Close()
}

// Err returns the error if any occurred during decoding.
func (s *eventStreamDecoder[E]) Error() error {
	return s.err
}
