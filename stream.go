package prompts

import (
	"iter"
)

// StreamDecoder is an interface for decoding a stream of events.
type StreamDecoder[E any] interface {
	// All returns an iterator over all events to be decoded.
	All() iter.Seq[E]
}

// StreamTransformer is a function that transforms an event into a different type.
type StreamTransformer[E any, T any] func(E) (T, error)

// Stream is the interface for a stream of events.
type Stream[T any] interface {
	// All returns an iterator over all events to be decoded.
	All() iter.Seq2[T, error]
	// Error returns the error if any occurred during decoding.
	Error() error
}

type stream[E any, T any] struct {
	decoder     StreamDecoder[E]
	transformer StreamTransformer[E, T]
	err         error
}

// NewStream creates a new stream from the given decoder and error.
func NewStream[E, T any](decoder StreamDecoder[E], transformer StreamTransformer[E, T]) Stream[T] {
	return &stream[E, T]{
		transformer: transformer,
		decoder:     decoder,
	}
}

// All returns an iterator over all events to be decoded.
func (s *stream[E, T]) All() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for e := range s.decoder.All() {
			if !yield(s.transformer(e)) {
				break
			}
		}
	}
}

// Error returns the error if any occurred during decoding.
func (s *stream[E, T]) Error() error {
	return s.err
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
