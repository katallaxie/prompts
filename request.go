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
	// RoleDeveloper is the developer role.
	RoleDeveloper Role = "developer"
	// RoleSystem is the system role.
	RoleSystem Role = "system"
	// RoleFunction is the function role.
	RoleFunction Role = "function"
	// RoleNone is the none role.
	RoleNone Role = ""
)

type ToolChoice string

const (
	// ToolChoiceAuto is the auto tool choice.
	ToolChoiceAuto ToolChoice = "auto"
	// ToolChoiceAll is the all tool choice.
	ToolChoiceNone ToolChoice = "none"
	// ToolChoiceRequired is the required tool choice.
	ToolChoiceRequired ToolChoice = "required"
)

type isChatCompletionTool interface {
	isChatCompletionTool()
}

// ChatCompletionTool represents a tool for the chat completion request.
type ChatCompletionTool struct {
	Tool isChatCompletionTool
}

func (c ChatCompletionTool) isChatCompletionTool() {}

// MarshalJSON marshals the chat completion tool into JSON.
func (c ChatCompletionTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Tool)
}

// ChatCompletionFuntionTool represents a function tool for the chat completion request.
type ChatCompletionFuntionTool struct {
	// Function is the function for the chat completion request.
	Function ChatCompletionFunctionDefintion `json:"function,omitempty"`
}

// MarshalJSON marshals the chat completion function tool into JSON.
func (c ChatCompletionFuntionTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name,omitempty"`
		// Description is the description of the function.
		Description string `json:"description,omitempty"`
		// Parameters is the parameters for the function.
		Parameters ChatCompletionFunctionParameters `json:"parameters,omitempty"`
		// Strict is a flag to indicate whether to strictly enforce the parameters.
		Strict bool `json:"strict,omitempty"`
	}{
		Type:        "function",
		Name:        c.Function.Name,
		Description: c.Function.Description,
		Parameters:  c.Function.Parameters,
		Strict:      c.Function.Strict,
	})
}

// ChatCompletionFunctionDefintion represents the function definition for the chat completion request.
type ChatCompletionFunctionDefintion struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is the description of the function.
	Description string `json:"description,omitempty"`
	// Parameters is the parameters for the function.
	Parameters ChatCompletionFunctionParameters `json:"parameters,omitempty"`
	// Strict is a flag to indicate whether to strictly enforce the parameters.
	Strict bool `json:"strict,omitempty"`
}

func (c ChatCompletionFuntionTool) isChatCompletionTool() {}

// ChatCompletionFunctionProperties represents the properties for the function tool.
type ChatCompletionFunctionProperties map[string]json.RawMessage

// ChatCompletionFunctionParameters represents the parameters for the function tool.
type ChatCompletionFunctionParameters struct {
	// Properties is the properties for the function tool.
	Properties ChatCompletionFunctionProperties `json:"properties,omitempty"`
	// Required is the required parameters for the function tool.
	Required []string `json:"required,omitempty"`
}

// MarshalJSON marshals the chat completion function parameters into JSON.
func (c ChatCompletionFunctionParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string                     `json:"type"`
		Properties map[string]json.RawMessage `json:"properties,omitempty"`
		Required   []string                   `json:"required,omitempty"`
	}{
		Type:       "object",
		Properties: c.Properties,
		Required:   c.Required,
	})
}

// ChatCompletionCustomTool represents a custom tool for the chat completion request.
type ChatCompletionCustomTool struct {
	// Custom is the custom tool for the chat completion request.
	Custom ChatCompletionCustomDefintion `json:"custom,omitempty"`
}

// MarshalJSON marshals the chat completion custom tool into JSON.
func (c ChatCompletionCustomTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string                        `json:"type"`
		Custom ChatCompletionCustomDefintion `json:"custom,omitempty"`
	}{
		Type:   "custom",
		Custom: c.Custom,
	})
}

// ChatCompletionCustomDefintion represents the custom definition for the chat completion request.
type ChatCompletionCustomDefintion struct {
	// Name is the name of the custom tool.
	Name string `json:"name"`
	// Description is the description of the custom tool.
	Description string `json:"description,omitempty"`
}

func (c ChatCompletionCustomTool) isChatCompletionTool() {}

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
		Type: "input_text",
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
	// Name is the name of the message sender (optional).
	Name string `json:"name,omitempty"`
}

// ChatCompletionRequest is the request for chat completion.
type ChatCompletionRequest struct {
	// Model is the model for the chat completion request.
	Model string `json:"model"`
	// Input is the list of messages for the chat completion request.
	Input []ChatCompletionMessage `json:"input"`
	// Instructions is the instructions for the chat completion request.
	Instructions string `json:"instructions,omitempty"`
	// Tools is the list of tools to use for the chat completion request.
	Tools []ChatCompletionTool `json:"tools,omitempty"`
	// ToolChoice is the tool choice for the chat completion request.
	ToolChoice ToolChoice `json:"tool_choice,omitempty"`
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

// WithInput sets the messages for the chat completion request.
func WithInput(msgs ...ChatCompletionMessage) Opt {
	return func(req *ChatCompletionRequest) {
		req.Input = msgs
	}
}

// WithInstructions sets the instructions for the chat completion request.
func WithInstructions(instructions string) Opt {
	return func(req *ChatCompletionRequest) {
		req.Instructions = instructions
	}
}

// WithTools sets the tools for the chat completion request.
func WithTools(tools ...ChatCompletionTool) Opt {
	return func(req *ChatCompletionRequest) {
		req.Tools = tools
	}
}
