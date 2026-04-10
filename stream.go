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
	// All returns all events.
	All() iter.Seq[E]
}

// Transformer is a function that transforms an event into a different type.
type Transformer[E any, T any] func(E) (T, error)

// Stream is the interface for a stream of events.
type Stream[E any, T any] interface {
	// All returns all events.
	All() iter.Seq2[T, error]
}

// Stream is a stream of events.
type stream[E any, T any] struct {
	decoder     Decoder[E]
	transformer Transformer[E, T]
}

// NewStream creates a new stream from the given decoder and error.
func NewStream[E, T any](decoder Decoder[E], transformer Transformer[E, T]) Stream[E, T] {
	return &stream[E, T]{
		transformer: transformer,
		decoder:     decoder,
	}
}

// All returns all events.
func (s *stream[E, T]) All() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for e := range s.decoder.All() {
			if !yield(s.transformer(e)) {
				break
			}
		}
	}
}

// Callbacks is an interfactor for a callback function.
func Callbacks[T any](v T, cb ...func(T) error) error {
	for _, c := range cb {
		if err := c(v); err != nil {
			return err
		}
	}

	return nil
}

// Events is a function that returns the events.
func Events[T any](stream iter.Seq2[T, error], cb ...func(T) error) error {
	for msg, err := range stream {
		if err != nil {
			return err
		}

		if err := Callbacks(msg, cb...); err != nil {
			return err
		}
	}

	return nil
}
