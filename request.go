package prompts

import (
	"fmt"
)

var _ fmt.Stringer = (*Role)(nil)

// Role is the role of the message sender.
type Role string

// String returns the string representation of the role.
func (r Role) String() string {
	return string(r)
}

// Available roles.
const (
	// RoleUser is the user role.
	RoleUser Role = "user"
	// RoleAssistant is the assistant role.
	RoleAssistant Role = "assistant"
	// RoleSystem is the system role.
	RoleSystem Role = "system"
	// RoleFunction is the function role.
	RoleFunction Role = "function"
	// RoleNone is the none role.
	RoleNone Role = ""
)

// ChatCompletionMessage is the message for chat completion.
type ChatCompletionMessage struct {
	// Role is the role of the message sender.
	Role Role `json:"role"`
	// Content is the content of the message.
	Content string `json:"content"`
	// Images is the list of images for the message.
	Images []Image `json:"images,omitempty"`
	// ToolCalls is the tool call for the message.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ChatCompletionRequest is the request for chat completion.
type ChatCompletionRequest struct {
	// Model is the model for the chat completion request.
	Model string `json:"model"`
	// Messages is the list of messages for the chat completion request.
	Messages []ChatCompletionMessage `json:"messages"`
	// MaxTokens is the maximum number of tokens for the chat completion request.
	MaxTokens *int `json:"max_tokens,omitzero"`
	// Temperature is the sampling temperature
	Temperature *float32 `json:"temperature,omitzero"`
	// Stream is a flag to enable streaming
	Stream bool `json:"stream,omitempty"`
	// TopP is the nucleus sampling parameter
	TopP *float64 `json:"top_p,omitzero"`
	// TopK is the number of top tokens to sample from
	TopK *int `json:"top_k,omitzero"`
	// Tools is the list of tools to use for the chat completion
	Tools []Tool `json:"tools,omitempty"`
}

// NewChatCompletionRequest creates a new chat completion request.
func NewChatCompletionRequest(msgs ...ChatCompletionMessage) *ChatCompletionRequest {
	return &ChatCompletionRequest{
		Messages: msgs,
	}
}

// NewStreamChatCompletionRequest creates a new chat completion request with streaming enabled.
func NewStreamChatCompletionRequest(msgs ...ChatCompletionMessage) *ChatCompletionRequest {
	return &ChatCompletionRequest{
		Messages: msgs,
		Stream:   true,
	}
}
