package prompts

import (
	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/pkg/slices"
	"github.com/katallaxie/pkg/utilx"
)

// Role is the role of the message sender
type Role string

// Available roles
const (
	// RoleUser is the user role
	RoleUser Role = "user"
	// RoleAssistant is the assistant role
	RoleAssistant Role = "assistant"
	// RoleSystem is the system role
	RoleSystem Role = "system"
	// RoleFunction is the function role
	RoleFunction Role = "function"
	// RoleNone is the none role
	RoleNone Role = ""
)

// ChatCompletionMessage is the message for chat completion
type ChatCompletionMessage struct {
	// Role is the role of the message sender
	Role Role `json:"role"`
	// Content is the content of the message
	Content string `json:"content"`
}

// ChatCompletionRequest is the request for chat completion
type ChatCompletionRequest struct {
	// Model is the model name
	Model string `json:"model"`
	// Messages is the list of messages
	Messages []ChatCompletionMessage `json:"messages"`
	// MaxTokens is the maximum number of tokens to generate
	MaxTokens *int `json:"max_tokens,omitzero"`
	// Temperature is the sampling temperature
	Temperature *float32 `json:"temperature,omitzero"`
	// Stream is a flag to enable streaming
	Stream *bool `json:"stream,omitempty"`
	// TopP is the nucleus sampling parameter
	TopP *float64 `json:"top_p,omitzero"`
	// TopK is the number of top tokens to sample from
	TopK *int `json:"top_k,omitzero"`
}

// SetModel sets the model for the chat completion request
func (r *ChatCompletionRequest) SetModel(model string) {
	r.Model = model
}

// AddMessage adds a message to the chat completion request
func (r *ChatCompletionRequest) AddMessages(msg ...ChatCompletionMessage) {
	r.Messages = slices.Append(r.Messages, msg...)
}

// SetMessages sets the messages for the chat completion request
func (r *ChatCompletionRequest) SetMessages(msg []ChatCompletionMessage) {
	r.Messages = msg
}

// SetMaxTokens sets the maximum number of tokens for the chat completion request
func (r *ChatCompletionRequest) SetMaxTokens(maxTokens int) {
	r.MaxTokens = cast.Ptr(maxTokens)
}

// SetTemperature sets the temperature for the chat completion request
func (r *ChatCompletionRequest) SetTemperature(temperature float32) {
	r.Temperature = cast.Ptr(temperature)
}

// SetTopP sets the top P for the chat completion request
func (r *ChatCompletionRequest) SetTopP(topP float64) {
	r.TopP = cast.Ptr(topP)
}

// SetTopK sets the top K for the chat completion request
func (r *ChatCompletionRequest) SetTopK(topK int) {
	r.TopK = cast.Ptr(topK)
}

// SetStream sets the stream flag for the chat completion request
func (r *ChatCompletionRequest) SetStream(stream bool) {
	r.Stream = cast.Ptr(stream)
}

// Index is the index for the chat completion
type Index struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionChoice is the choice for chat completion
type ChatCompletionChoice struct {
	// Message is the message for the choice
	Message Index `json:"message,omitempty"`
	// FinishReason is the reason for finishing
	FinishReason string `json:"finish_reason,omitempty"`
	// Delta is the delta for the choice
	Delta Index `json:"delta,omitempty"`
	// Index is the index for the choice
	Index uint `json:"index,omitempty"`
}

// String returns the string representation of the choice
func (c ChatCompletionChoice) String() string {
	if utilx.NotEmpty(c.Delta.Content) {
		return c.Delta.Content
	}

	return c.Message.Content
}
