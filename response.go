package prompts

import (
	"fmt"
	"io"
	"strings"

	"github.com/katallaxie/pkg/utilx"
)

// FinishReason is the reason for finishing a chat completion.
type FinishReason string

const (
	// FinishReasonStop indicates that the chat completion was finished because the model stopped generating content.
	FinishReasonStop FinishReason = "stop"
	// FinishReasonLength indicates that the chat completion was finished because the maximum length was reached.
	FinishReasonLength FinishReason = "length"
	// FinishReasonContentFilter indicates that the chat completion was finished because the content filter was triggered.
	FinishReasonContentFilter FinishReason = "content_filter"
	// FinishReasonUnknown indicates that the chat completion was finished for an unknown reason.
	FinishReasonUnknown FinishReason = ""
)

var _ fmt.Stringer = (*FinishReason)(nil)

// String returns the string representation of the finish reason.
func (f FinishReason) String() string {
	return string(f)
}

// ChatCompletionChoiceIndex is the index for the chat completion.
type ChatCompletionChoiceIndex struct {
	// Role is the role of the message sender.
	Role Role `json:"role"`
	// Content is the content of the message.
	Content string `json:"content"`
}

// CompletionUsage represents the usage of the chat completion.
type CompletionUsage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens,omitempty"`
	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens,omitempty"`
	// TotalTokens is the total number of tokens used in the chat completion.
	TotalTokens int `json:"total_tokens,omitempty"`
}

var _ fmt.Stringer = (*ChatCompletionResponse)(nil)

// ChatCompletionChoice is the choice for chat completion.
type ChatCompletionChoice struct {
	// Message is the message for the choice
	Message ChatCompletionChoiceIndex `json:"message,omitempty"`
	// FinishReason is the reason for finishing
	FinishReason FinishReason `json:"finish_reason,omitempty"`
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
	// Citations is the list of citations returned in the response
	Citations []string `json:"citations,omitempty"`
	// SearchResults is the list of search results returned in the response
	SearchResults []SearchResult `json:"search_results,omitempty"`
	// SystemFingerprint is the system fingerprint returned in the response
	SystemFingerprint string `json:"system_fingerprint,omitempty"`
	// CompletionUsage is the usage of the chat completion returned in the response
	CompletionUsage CompletionUsage `json:"usage,omitempty"`
}

// SearchResult represents a search result structure for chat completion API.
type SearchResult struct {
	// Title is the title of the search result
	Title string `json:"title,omitempty"`
	// URL is the URL of the search result
	URL string `json:"url,omitempty"`
	// Snippet is the snippet of the search result
	Snippet string `json:"snippet,omitempty"`
	// Source is the source of the search result
	Source string `json:"source,omitempty"`
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
