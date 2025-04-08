package prompts

import (
	"iter"
)

// Decoder is an interface for decoding SSE streams.
type Decoder[E any] interface {
	// Close closes the decoder.
	Close() error
	// Error returns the error if any occurred during decoding.
	Error() error
	// Next returns true if there is a next event.
	Next() iter.Seq[E]
}

// Transformer is a function that transforms an event into a different type.
type Transformer[E any, T any] func(E) (T, error)

// Stream is a stream of events.
type Stream[E any, T any] struct {
	decoder     Decoder[E]
	transformer Transformer[E, T]
	err         error
}

// NewStream creates a new stream from the given decoder and error.
func NewStream[E any, T any](decoder Decoder[E], transformer Transformer[E, T]) *Stream[E, T] {
	return &Stream[E, T]{
		transformer: transformer,
		decoder:     decoder,
	}
}

// Scan returns a sequence of events.
func (s *Stream[E, T]) Next() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := range s.decoder.Next() {
			msg, err := s.transformer(e)
			if err != nil {
				s.err = err
				break
			}

			if !yield(msg) {
				break
			}
		}

		if s.decoder.Error() != nil {
			s.err = s.decoder.Error()
		}
	}
}

// Err returns the error if any occurred during streaming.
func (s *Stream[E, T]) Error() error {
	return s.err
}

// Done returns true if the stream has ended.
func (s *Stream[E, T]) Close() error {
	if s.decoder == nil {
		return nil
	}

	return s.decoder.Close()
}
