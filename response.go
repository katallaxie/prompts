package prompts

import (
	"fmt"
	"io"
	"strings"

	"github.com/katallaxie/pkg/utilx"
)

// ChatCompletionChoiceIndex is the index for the chat completion.
type ChatCompletionChoiceIndex struct {
	// Role is the role of the message sender.
	Role Role `json:"role"`
	// Content is the content of the message.
	Content string `json:"content"`
}

var _ fmt.Stringer = (*ChatCompletionResponse)(nil)

// ChatCompletionChoice is the choice for chat completion.
type ChatCompletionChoice struct {
	// Message is the message for the choice
	Message ChatCompletionChoiceIndex `json:"message,omitempty"`
	// FinishReason is the reason for finishing
	FinishReason string `json:"finish_reason,omitempty"`
	// Delta is the delta for the choice
	Delta ChatCompletionChoiceIndex `json:"delta,omitempty"`
	// Index is the index for the choice
	Index uint `json:"index,omitempty"`
}

// String returns the string representation of the choice.
func (c ChatCompletionChoice) String() string {
	if utilx.NotEmpty(c.Delta.Content) {
		return c.Delta.Content
	}

	return c.Message.Content
}

var (
	_ fmt.Stringer = (*ChatCompletionResponse)(nil)
	_ io.Reader    = (*ChatCompletionResponse)(nil)
)

// NewChatCompletionResponse creates a new ChatCompletionResponse with the given choices.
func NewChatCompletionResponse(choices ...ChatCompletionChoice) *ChatCompletionResponse {
	return &ChatCompletionResponse{
		Choices: choices,
	}
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	// ID is the unique identifier for the response
	ID string `json:"id,omitempty"`
	// Object is the type of object returned
	Object string `json:"object,omitempty"`
	// Created is the timestamp of when the response was created
	Created int64 `json:"created,omitempty"`
	// Model is the model used for the response
	Model string `json:"model"`
	// Choices is the list of choices returned in the response
	Choices []ChatCompletionChoice `json:"choices"`
}

// String returns the string representation of the response.
func (c *ChatCompletionResponse) String() string {
	var out strings.Builder

	for _, choice := range c.Choices {
		out.WriteString(choice.String())
	}

	return out.String()
}

// Read implements the io.Reader interface for ChatCompletionResponse.
// It reads the content of the choices and writes it to the provided byte slice.
func (c *ChatCompletionResponse) Read(p []byte) (n int, err error) {
	var out strings.Builder

	for _, choice := range c.Choices {
		out.WriteString(choice.String())
	}

	copy(p, out.String())

	return len(out.String()), io.EOF
}
