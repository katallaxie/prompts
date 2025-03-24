package prompts

import (
	"context"
	"errors"
)

// Model is a chat model.
type Model string

// String returns the string representation of the chat model.
func (m Model) String() string {
	return string(m)
}

// Role is a chat role.
type Role string

// String returns the string representation of the chat model.
func (r Role) String() string {
	return string(r)
}

const (
	// RoleUser indicates that a message was send by a user.
	RoleUser Role = "user"
	// RoleHuman indicates that a message was send by a human.
	RoleHuman Role = "human"
	// RoleAI indicates that a message was send by an AI.
	RoleAI Role = "ai"
	// RoleSystem indicates that a message was send by the system.
	RoleSystem Role = "system"
	// RoleAssistant indicates that a message was send by an assistant.
	RoleAssistant Role = "assistant"
	// RoleTool indicates that a message was send by a tool.
	RoleTool Role = "tool"
	// RoleFunction indicates that a message was send by a function.
	RoleFunction Role = "function"
	// RoleGeneric indicates that a message was send by a generic role.
	RoleGeneric Role = "generic"
)

var (
	ErrUnsupportedRole  = errors.New("chats: unsupported role")
	ErrUnsupportedModel = errors.New("chats: unsupported model")
)

// Message is a prompt message.
type Message interface {
	// GetRole returns the role of the message.
	GetRole() Role
	// GetContent returns the content of the message.
	GetContent() string
}

var (
	_ Message = (*UserMessage)(nil)
	_ Message = (*HumanMessage)(nil)
	_ Message = (*AIMessage)(nil)
	_ Message = (*SystemMessage)(nil)
	_ Message = (*GenericMessage)(nil)
	_ Message = (*ToolMessage)(nil)
)

// Human chat message.
type HumanMessage struct {
	Content string `json:"content,omitempty"`
}

// GetContent returns the content of the message.
func (m *HumanMessage) GetContent() string {
	return m.Content
}

// GetRole returns the role of the message.
func (m *HumanMessage) GetRole() Role {
	return RoleHuman
}

// UserMessage is a user chat message.
type UserMessage struct {
	Content string `json:"content,omitempty"`
}

// GetContent returns the content of the message.
func (m *UserMessage) GetContent() string {
	return m.Content
}

// GetRole returns the role of the message.
func (m *UserMessage) GetRole() Role {
	return RoleUser
}

// AIMessage is an AI chat message.
type AIMessage struct {
	Content string `json:"content,omitempty"`
}

// GetContent returns the content of the message.
func (m *AIMessage) GetContent() string {
	return m.Content
}

// GetRole returns the role of the message.
func (m *AIMessage) GetRole() Role {
	return RoleAI
}

// SystemMessage is a system chat message.
type SystemMessage struct {
	Content string `json:"content,omitempty"`
}

// GetContent returns the content of the message.
func (m *SystemMessage) GetContent() string {
	return m.Content
}

// GetRole returns the role of the message.
func (m *SystemMessage) GetRole() Role {
	return RoleAI
}

// GenericMessage is a generic chat message.
type GenericMessage struct {
	Content string `json:"content,omitempty"`
	Role    Role   `json:"role,omitempty"`
	Name    string `json:"name,omitempty"`
}

// GetContent returns the content of the message.
func (m *GenericMessage) GetContent() string {
	return m.Content
}

// GetRole returns the role of the message.
func (m *GenericMessage) GetRole() Role {
	return m.Role
}

// GetName returns the name of the message.
func (m *GenericMessage) GetName() string {
	return m.Name
}

// ToolMessage is a tool chat message.
type ToolMessage struct {
	// ID is the ID of the message.
	ID string `json:"tool_call_id,omitempty"`
	// Content is the content of the message.
	Content string `json:"content,omitempty"`
}

// GetContent returns the content of the message.
func (m *ToolMessage) GetContent() string {
	return m.Content
}

// GetRole returns the role of the message.
func (m *ToolMessage) GetRole() Role {
	return RoleTool
}

// GetID returns the ID of the message.
func (m *ToolMessage) GetID() string {
	return m.ID
}

// Prompt is a prompt.
type Prompt struct {
	// Model is the model to use.
	Model Model `json:"model"`
	// Messages are the messages to use.
	Messages []Message `json:"messages"`
	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens int `json:"max_tokens,omitzero"`
	// Temperature is the temperature.
	Temperature *float64 `json:"temperature,omitempty"`
	// TopP is the top p.
	TopP *float64 `json:"top_p,omitempty"`
	// TopK is the top k.
	TopK *uint `json:"top_k,omitempty"`
}

// CompletionChoice is a completion choice.
type CompletionChoice struct {
	// Message is the message of the choice.
	Message Message `json:"message,omitempty"`
	// FinishReason is the finish reason of the choice.
	FinishReason string `json:"finish_reason,omitempty"`
	// Delete is the delete of the choice.
	Delete Message `json:"delete,omitempty"`
	// Index is the index of the choice.
	Index uint `json:"index,omitempty"`
}

// Completion is a completion.
type Completion struct {
	// ID is the ID of the response.
	ID string `json:"id,omitempty"`
	// Object is the object of the response.
	Object string `json:"object,omitempty"`
	// Created is the creation time of the response.
	Created int64 `json:"created,omitempty"`
	// Model is the model of the response.
	Model Model `json:"model"`
	// Choices are the choices of the response.
	Choices []CompletionChoice `json:"choices"`
}

// NewCompletion returns a new completion.
func NewCompletion(out chan CompletionChoice) *Completion {
	c := new(Completion)

	return c
}

// Promptable is a promptable.
type Promptable interface {
	// Prompt prompts a completion.
	Complete(ctx context.Context, prompt *Prompt) (chan any, error)
}
