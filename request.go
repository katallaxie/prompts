package prompts

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

// DefaultClient is the default HTTP client for the chat completion request.
var DefaultClient = &http.Client{
	Timeout: defaultTimeout,
}

// Defaults returns the default options for the chat completion request.
func Defaults() *Opts {
	return &Opts{
		Client: DefaultClient,
	}
}

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

// ToolCall represents a tool call that can be used in a chat completion response.
type ToolCall struct {
	ToolCall isToolCall
}

// Reset resets the tool call for the chat completion response.
func (t *ToolCall) Reset() {
	*t = ToolCall{}
}

type isToolCall interface {
	isToolCall()
}

// NewToolCall creates a new tool call for the chat completion response.
func NewToolCall() ToolCall {
	return ToolCall{}
}

// ToolCallFunction represents a tool call function that can be used in a chat completion response.
type ToolCallFunction struct {
	// ID is the ID of the function to call when the tool is used.
	ID string `json:"id"`
	// CallID is the call ID for the function call.
	CallID string `json:"call_id"`
	// Name is the name of the function to call when the tool is used.
	Name string `json:"name"`
	// Arguments is the arguments to pass to the function when the tool is used.
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// MarshalJSON marshals the tool call function into JSON.
func (t ToolCallFunction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string                 `json:"type"`
		ID        string                 `json:"id"`
		CallID    string                 `json:"call_id"`
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments,omitempty"`
	}{
		Type:      "function",
		ID:        t.ID,
		CallID:    t.CallID,
		Name:      t.Name,
		Arguments: t.Arguments,
	})
}

func (t ToolCallFunction) isToolCall() {}

// Tool represents the tools for the chat completion request.
type Tool struct {
	Tool isTool
}

type isTool interface {
	isTool()
}

// Reset resets the tools for the chat completion request.
func (t *Tool) Reset() {
	*t = Tool{}
}

// NewTool creates a new tool for the chat completion request.
func NewTool() Tool {
	return Tool{}
}

// ToolFunction represents a function to call when a tool is used.
type ToolFunction struct {
	// Name is the name of the function to call when the tool is used.
	Name string `json:"name"`
	// Description is the description of the function to call when the tool is used.
	Description string `json:"description,omitempty"`
	// Strict indicates whether the function call should be strictly enforced.
	Strict bool `json:"strict,omitempty"`
	// Parameters is the parameters to pass to the function when the tool is used.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	// Required is the list of required parameters for the function call.
	Required []string `json:"required,omitempty"`
	// AdditionalProperties indicates whether additional properties are allowed for the function call.
	AdditionalProperties bool `json:"additional_properties,omitempty"`
}

// MarshalJSON marshals the tool call function into JSON.
func (t ToolFunction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type                 string                 `json:"type"`
		Name                 string                 `json:"name"`
		Description          string                 `json:"description,omitempty"`
		Strict               bool                   `json:"strict,omitempty"`
		Parameters           map[string]interface{} `json:"parameters,omitempty"`
		Required             []string               `json:"required,omitempty"`
		AdditionalProperties bool                   `json:"additionalProperties,omitempty"` //nolint:tagliatelle
	}{
		Type:                 "function",
		Name:                 t.Name,
		Description:          t.Description,
		Strict:               t.Strict,
		Parameters:           t.Parameters,
		Required:             t.Required,
		AdditionalProperties: t.AdditionalProperties,
	})
}

func (t ToolFunction) isToolCall() {}

// ChatCompletionMessageContent is the content of a chat completion message.
type ChatCompletionMessageContent struct {
	Content isChatCompletionMessageContent
}

// MarhalJSON marshals the chat completion message content into JSON.
func (c ChatCompletionMessageContent) MarshalJSON() ([]byte, error) {
	if text, ok := c.GetText(); ok {
		return json.Marshal(text)
	}

	return json.Marshal(nil) // Return null if the content is not text
}

type isChatCompletionMessageContent interface {
	isChatCompletionMessageContent()
}

// NewChatCompletionMessageContent creates a new chat completion message content.
func NewChatCompletionMessageContent() ChatCompletionMessageContent {
	return ChatCompletionMessageContent{}
}

// Reset resets the chat completion message content.
func (c *ChatCompletionMessageContent) Reset() {
	*c = ChatCompletionMessageContent{}
}

// GetText returns the text content of the chat completion message content.
func (c ChatCompletionMessageContent) GetText() (ChatCompletionMessageContentText, bool) {
	if text, ok := c.Content.(ChatCompletionMessageContentText); ok {
		return text, true
	}

	return ChatCompletionMessageContentText{}, false
}

// GetImage returns the image content of the chat completion message content.
func (c ChatCompletionMessageContent) GetImage() (ChatCompletionMessageContentImage, bool) {
	if image, ok := c.Content.(ChatCompletionMessageContentImage); ok {
		return image, true
	}

	return ChatCompletionMessageContentImage{}, false
}

// GetFile returns the file content of the chat completion message content.
func (c ChatCompletionMessageContent) GetFile() (ChatCompletionMessageContentFile, bool) {
	if file, ok := c.Content.(ChatCompletionMessageContentFile); ok {
		return file, true
	}

	return ChatCompletionMessageContentFile{}, false
}

// ChatCompletionMessageContentText is the text content of a chat completion message.
type ChatCompletionMessageContentText struct {
	// Text is the text of the content.
	Text string `json:"text"`
}

// MarshalJSON marshals the chat completion message content text into JSON.
func (c ChatCompletionMessageContentText) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{
		Type: "text",
		Text: c.Text,
	})
}

func (c ChatCompletionMessageContentText) isChatCompletionMessageContent() {}

// ChatCompletionMessageContentImage is the image content of a chat completion message.
type ChatCompletionMessageContentImage struct {
	Image Image `json:"image"`
}

func (c ChatCompletionMessageContentImage) isChatCompletionMessageContent() {}

// ChatCompletionMessageContentFile is the file content of a chat completion message.
type ChatCompletionMessageContentFile struct {
	File File `json:"file"`
}

func (c ChatCompletionMessageContentFile) isChatCompletionMessageContent() {}

// File is the file for the chat completion message content.
type File struct {
	// Name is the name of the file.
	Name string `json:"name"`
	// URL is the URL of the file.
	URL string `json:"url"`
}

// Image is a type that represents an image.
type Image struct {
	// URL is the URL of the image.
	URL string `json:"url,omitempty"`
	// Base64 is the base64 encoding of the image.
	Base64 string `json:"base64,omitempty"`
	// Name is the name of the image.
	Name string `json:"name,omitempty"`
}

// Encode encodes the image into a string.
func (i Image) Encode(data []byte) string {
	i.Base64 = base64.StdEncoding.EncodeToString(data)
	return i.Base64
}

// NewImage creates a new image from the given data.
func NewImage(data []byte) Image {
	var img Image
	img.Encode(data)
	return img
}

// ChatCompletionMessage is the message for chat completion.
type ChatCompletionMessage struct {
	// Role is the role of the message sender.
	Role Role `json:"role"`
	// Content is the content of the message.
	Content []ChatCompletionMessageContent `json:"content"`
	// ToolCalls is the tool call for the message.
	// ToolCalls []ToolCall `json:"tool_calls,omitempty"`
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
	// Opts is the options for the chat completion request.
	Opts *Opts `json:"-"`
}

// NewChatCompletionRequest creates a new chat completion request.
func NewChatCompletionRequest(opts ...Opt) *ChatCompletionRequest {
	req := new(ChatCompletionRequest)
	req.Opts = Defaults()

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// NewStreamChatCompletionRequest creates a new chat completion request with streaming enabled.
func NewStreamChatCompletionRequest(opts ...Opt) *ChatCompletionRequest {
	req := new(ChatCompletionRequest)
	req.Stream = true
	req.Opts = Defaults()

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// Opt is a function that configures the options for the Prompter.
type Opt func(*ChatCompletionRequest)

// Opts is the options for configuring the Prompter.
type Opts struct {
	// BaseURL is the base URL.
	BaseURL string `json:"base_url"`
	// ApiKey is the API key.
	ApiKey string `json:"api_key"`
	// Headers are the headers to include in the request.
	Headers map[string]string `json:"headers"`
	// Client is the HTTP client.
	Client *http.Client `json:"-"`
}

// WithURL configures the base URL.
func WithURL(url string) Opt {
	return func(o *ChatCompletionRequest) {
		o.Opts.BaseURL = url
	}
}

// WithApiKey configures the API key.
func WithApiKey(apiKey string) Opt {
	return func(o *ChatCompletionRequest) {
		o.Opts.ApiKey = apiKey
	}
}

// WithClient configures the HTTP client.
func WithClient(client *http.Client) Opt {
	return func(o *ChatCompletionRequest) {
		o.Opts.Client = client
	}
}

// WithBaseURL configures the base URL.
func WithBaseURL(url string) Opt {
	return func(o *ChatCompletionRequest) {
		o.Opts.BaseURL = url
	}
}

// WithMessages sets the messages for the chat completion request.
func WithMessages(msgs ...ChatCompletionMessage) Opt {
	return func(req *ChatCompletionRequest) {
		req.Messages = msgs
	}
}
