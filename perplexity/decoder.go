package perplexity

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"iter"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/pkg/utilx"
	"github.com/katallaxie/prompts"
)

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

// NewDecoder creates a new Decoder based on the content type of the response.
func NewDecoder(body io.ReadCloser) prompts.StreamDecoder[Event] {
	if utilx.Empty(body) {
		return nil
	}

	scanner := bufio.NewScanner(body)
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)

	return &eventStreamDecoder{rc: body, scn: scanner}
}

type eventStreamDecoder struct {
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

// All returns an iterator over all events to be decoded.
func (s *eventStreamDecoder) All() iter.Seq[Event] {
	return func(yield func(Event) bool) {
		for s.scn.Scan() {
			event := cast.Zero[Event]()

			b := s.scn.Bytes()

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

		if s.scn.Err() != nil {
			s.err = s.scn.Err()
		}
	}
}

// Err returns the error if any occurred during decoding.
func (s *eventStreamDecoder) Error() error {
	return s.err
}
