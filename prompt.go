package prompts

import (
	"context"
)

// // Generator is a type that represents a generator of chat completion responses.
// // It is an iterator that yields chat completion responses or errors.
// type Generator iter.Seq2[*Response, error]

// // Decoder is a function type that implements the Decoder interface.
// type Decoder[E any] interface {
// 	Decode(io.ReadCloser) iter.Seq[E]
// }

// // Transformer is a function type that implements the Transformer interface.
// type Transformer[E any] interface {
// 	Transform(iter.Seq[E]) Generator
// }

// Responder is the interface for sending a chat completion request and receiving a response.
type Responder[I, O any] interface {
	// Respond sends a chat completion request and returns the response.
	Respond(ctx context.Context, in I) (O, error)
}

// Prompter is the interface for sending a chat completion request and receiving a response.
type Prompter[I, O any] interface {
	Responder[I, O]
}
