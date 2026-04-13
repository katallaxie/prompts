package prompts

import (
	"io"
	"iter"
)

// StreamDecoder is a function type that implements the StreamDecoder interface.
type StreamDecoder[E any] func(io.ReadCloser) iter.Seq[E]

// StreamTransformer is a function that transforms an event into a different type.
type StreamTransformer[E any] func(E) (*ChatCompletionResponse, error)

// StreamIterator is an interface for iterating over a stream of events.
type StreamIterator func(iter.Seq2[*ChatCompletionResponse, error]) error

// Stream is the interface for a stream of events.
type Stream[T any] interface {
	// All returns an iterator over all events to be decoded.
	All() iter.Seq2[T, error]
	// Error returns the error if any occurred during decoding.
	Error() error
}

type stream[E any] struct {
	body        io.ReadCloser
	decoder     StreamDecoder[E]
	transformer StreamTransformer[E]
	err         error
}

// NewStream creates a new stream from the given decoder and error.
func NewStream[E any](body io.ReadCloser, decoder StreamDecoder[E], transformer StreamTransformer[E]) Stream[*ChatCompletionResponse] {
	return &stream[E]{
		transformer: transformer,
		body:        body,
		decoder:     decoder,
	}
}

// All returns an iterator over all events to be decoded.
func (s *stream[E]) All() iter.Seq2[*ChatCompletionResponse, error] {
	return func(yield func(*ChatCompletionResponse, error) bool) {
		for e := range s.decoder(s.body) {
			if !yield(s.transformer(e)) {
				break
			}
		}
	}
}

// Error returns the error if any occurred during decoding.
func (s *stream[E]) Error() error {
	return s.err
}
