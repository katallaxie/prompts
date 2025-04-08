package prompts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/katallaxie/pkg/utilx"
	"github.com/tidwall/gjson"
)

// Decoder is an interface for decoding SSE streams.
type Decoder interface {
	Close() error
	Error() error
	Event() Event
	Next() bool
}

// NewDecoder creates a new Decoder based on the content type of the response.
func NewDecoder(res *http.Response) Decoder {
	if utilx.Or(utilx.Empty(res), utilx.Empty(res.Body)) {
		return nil
	}

	contentType := res.Header.Get("content-type")
	if read, ok := decoderTypes[contentType]; ok {
		return read(res.Body)
	}

	scanner := bufio.NewScanner(res.Body)
	return &eventStreamDecoder{rc: res.Body, scn: scanner}
}

var decoderTypes = map[string](func(io.ReadCloser) Decoder){}

// RegisterDecoder registers a new decoder for a specific content type.
func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	decoderTypes[strings.ToLower(contentType)] = decoder
}

// Event represents an event in the SSE stream.
type Event struct {
	// Type is the type of the event.
	Type string
	// Data is the data of the event.
	Data []byte
}

var _ Decoder = (*eventStreamDecoder)(nil)

type eventStreamDecoder struct {
	evt Event
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

// NewEventStreamDecoder creates a new event stream decoder.
// nolint:gocyclo
func (s *eventStreamDecoder) Next() bool {
	if utilx.NotEmpty(s.err) {
		return false
	}

	event := ""
	data := bytes.NewBuffer(nil)

	for s.scn.Scan() {
		fmt.Print("new message")
		txt := s.scn.Bytes()

		if len(txt) == 0 {
			s.evt = Event{
				Type: event,
				Data: data.Bytes(),
			}

			return true
		}

		name, value, _ := bytes.Cut(txt, []byte(":"))

		if len(value) > 0 && value[0] == ' ' {
			value = value[1:]
		}

		switch string(name) {
		case "":
			continue
		case "event":
			event = string(value)
		case "data":
			_, s.err = data.Write(value)
			if s.err != nil {
				break
			}
			s.err = data.WriteByte('\n')
		}
	}

	if s.scn.Err() != nil {
		s.err = s.scn.Err()
	}

	return false
}

// Event returns the current event.
func (s *eventStreamDecoder) Event() Event {
	return s.evt
}

// Close closes the decoder.
func (s *eventStreamDecoder) Close() error {
	return s.rc.Close()
}

// Err returns the error if any occurred during decoding.
func (s *eventStreamDecoder) Error() error {
	return s.err
}

// Stream is a stream of events.
type Stream[T any] struct {
	decoder Decoder
	cur     T
	err     error
	done    bool
}

// NewStream creates a new stream from the given decoder and error.
func NewStream[T any](decoder Decoder, err error) *Stream[T] {
	return &Stream[T]{
		decoder: decoder,
		err:     err,
	}
}

// Next returns false if the stream has ended or an error occurred.
// Call Stream.Current() to get the current value.
// Call Stream.Err() to get the error.
//
//		for stream.Next() {
//			data := stream.Current()
//		}
//
//	 	if stream.Err() != nil {
//			...
//	 	}
//
// nolint:gocyclo
func (s *Stream[T]) Next() bool {
	if s.err != nil {
		return false
	}

	for s.decoder.Next() {
		if s.done {
			continue
		}

		if bytes.HasPrefix(s.decoder.Event().Data, []byte("[DONE]")) {
			// In this case we don't break because we still want to iterate through the full stream.
			s.done = true
			continue
		}

		var nxt T
		if s.decoder.Event().Type == "" || strings.HasPrefix(s.decoder.Event().Type, "response.") {
			ep := gjson.GetBytes(s.decoder.Event().Data, "error")
			if ep.Exists() {
				s.err = fmt.Errorf("received error while streaming: %s", ep.String())
				return false
			}
			s.err = json.Unmarshal(s.decoder.Event().Data, &nxt)
			if s.err != nil {
				return false
			}
			s.cur = nxt
			return true
		} else {
			ep := gjson.GetBytes(s.decoder.Event().Data, "error")
			if ep.Exists() {
				s.err = fmt.Errorf("received error while streaming: %s", ep.String())
				return false
			}
			event := s.decoder.Event().Type
			data := s.decoder.Event().Data
			s.err = json.Unmarshal([]byte(fmt.Sprintf(`{ "event": %q, "data": %s }`, event, data)), &nxt)
			if s.err != nil {
				return false
			}
			s.cur = nxt
			return true
		}
	}

	// decoder.Next() may be false because of an error
	s.err = s.decoder.Error()

	return false
}

// Current returns the current value of the stream.
func (s *Stream[T]) Current() T {
	return s.cur
}

// Err returns the error if any occurred during streaming.
func (s *Stream[T]) Error() error {
	return s.err
}

// Done returns true if the stream has ended.
func (s *Stream[T]) Close() error {
	if s.decoder == nil {
		// already closed
		return nil
	}
	return s.decoder.Close()
}
