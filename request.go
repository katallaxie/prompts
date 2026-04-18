package prompts

import (
	"encoding/base64"
	"encoding/json"
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

type isResponseTool interface {
	isResponseTool()
}

// ResponseTool represents a tool for the chat completion request.
type ResponseTool struct {
	Tool isResponseTool
}

func (c ResponseTool) isResponseTool() {}

// MarshalJSON marshals the response tool into JSON.
func (c ResponseTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Tool)
}

// ResponseFunctionTool represents a function tool for the chat completion request.
type ResponseFunctionTool struct {
	// Function is the function for the chat completion request.
	Function ResponseFunctionDefinition `json:"function,omitempty"`
}

// MarshalJSON marshals the response function tool into JSON.
func (c ResponseFunctionTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name,omitempty"`
		// Description is the description of the function.
		Description string `json:"description,omitempty"`
		// Parameters is the parameters for the function.
		Parameters ResponseFunctionParameters `json:"parameters,omitempty"`
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

// ResponseFunctionDefinition represents the function definition for the chat completion request.
type ResponseFunctionDefinition struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is the description of the function.
	Description string `json:"description,omitempty"`
	// Parameters is the parameters for the function.
	Parameters ResponseFunctionParameters `json:"parameters,omitempty"`
	// Strict is a flag to indicate whether to strictly enforce the parameters.
	Strict bool `json:"strict,omitempty"`
}

func (c ResponseFunctionTool) isResponseTool() {}

// ResponseFunctionProperties represents the properties for the function tool.
type ResponseFunctionProperties map[string]json.RawMessage

// ResponseFunctionParameters represents the parameters for the function tool.
type ResponseFunctionParameters struct {
	// Properties is the properties for the function tool.
	Properties ResponseFunctionProperties `json:"properties,omitempty"`
	// Required is the required parameters for the function tool.
	Required []string `json:"required,omitempty"`
}

// MarshalJSON marshals the response function parameters into JSON.
func (c ResponseFunctionParameters) MarshalJSON() ([]byte, error) {
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

// ResponseCustomTool represents a custom tool for the chat completion request.
type ResponseCustomTool struct {
	// Custom is the custom tool for the chat completion request.
	Custom ResponseCustomDefinition `json:"custom,omitempty"`
}

// MarshalJSON marshals the response custom tool into JSON.
func (c ResponseCustomTool) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string                   `json:"type"`
		Custom ResponseCustomDefinition `json:"custom,omitempty"`
	}{
		Type:   "custom",
		Custom: c.Custom,
	})
}

// ResponseCustomDefinition represents the custom definition for the chat completion request.
type ResponseCustomDefinition struct {
	// Name is the name of the custom tool.
	Name string `json:"name"`
	// Description is the description of the custom tool.
	Description string `json:"description,omitempty"`
}

func (c ResponseCustomTool) isResponseTool() {}

// ResponseMessageContent is the content of a response message.
type ResponseMessageContent struct {
	Content isResponseMessageContent
}

// MarshalJSON marshals the response message content into JSON.
func (c ResponseMessageContent) MarshalJSON() ([]byte, error) {
	if text, ok := c.GetText(); ok {
		return json.Marshal(text)
	}

	return json.Marshal(nil) // Return null if the content is not text
}

type isResponseMessageContent interface {
	isResponseMessageContent()
}

// NewResponseMessageContent creates a new response message content.
func NewResponseMessageContent() ResponseMessageContent {
	return ResponseMessageContent{}
}

// Reset resets the response message content.
func (c *ResponseMessageContent) Reset() {
	*c = ResponseMessageContent{}
}

// GetText returns the text content of the response message content.
func (c ResponseMessageContent) GetText() (ResponseMessageContentText, bool) {
	if text, ok := c.Content.(ResponseMessageContentText); ok {
		return text, true
	}

	return ResponseMessageContentText{}, false
}

// GetImage returns the image content of the response message content.
func (c ResponseMessageContent) GetImage() (ResponseMessageContentImage, bool) {
	if image, ok := c.Content.(ResponseMessageContentImage); ok {
		return image, true
	}

	return ResponseMessageContentImage{}, false
}

// GetFile returns the file content of the response message content.
func (c ResponseMessageContent) GetFile() (ResponseMessageContentFile, bool) {
	if file, ok := c.Content.(ResponseMessageContentFile); ok {
		return file, true
	}

	return ResponseMessageContentFile{}, false
}

// ResponseMessageContentText is the text content of a response message.
type ResponseMessageContentText struct {
	// Text is the text of the content.
	Text string `json:"text"`
}

// MarshalJSON marshals the response message content text into JSON.
func (c ResponseMessageContentText) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{
		Type: "input_text",
		Text: c.Text,
	})
}

func (c ResponseMessageContentText) isResponseMessageContent() {}

// ResponseMessageContentImage is the image content of a response message.
type ResponseMessageContentImage struct {
	Image Image `json:"image"`
}

func (c ResponseMessageContentImage) isResponseMessageContent() {}

// ResponseMessageContentFile is the file content of a response message.
type ResponseMessageContentFile struct {
	File File `json:"file"`
}

func (c ResponseMessageContentFile) isResponseMessageContent() {}

// File is the file for the response message content.
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

// ResponseInput is the message for chat completion.
type ResponseInput struct {
	// Role is the role of the message sender.
	Role Role `json:"role"`
	// Content is the content of the message.
	Content []ResponseMessageContent `json:"content"`
	// Name is the name of the message sender (optional).
	Name string `json:"name,omitempty"`
}

// ResponseRequest is the request for chat completion.
type ResponseRequest struct {
	// Model is the model for the chat completion request.
	Model string `json:"model"`
	// Input is the list of messages for the chat completion request.
	Input []ResponseInput `json:"input"`
	// Instructions is the instructions for the chat completion request.
	Instructions string `json:"instructions,omitempty"`
	// Tools is the list of tools to use for the chat completion request.
	Tools []ResponseTool `json:"tools,omitempty"`
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
}

// RequestOpt is a function type for configuring the ResponseRequest.
type RequestOpt func(*ResponseRequest)

// NewResponseRequest creates a new chat completion request with the given options.
func NewResponseRequest(opts ...RequestOpt) *ResponseRequest {
	req := new(ResponseRequest)

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// WithInput sets the messages for the chat completion request.
func WithInput(msgs ...ResponseInput) RequestOpt {
	return func(req *ResponseRequest) {
		req.Input = msgs
	}
}

// WithInstructions sets the instructions for the chat completion request.
func WithInstructions(instructions string) RequestOpt {
	return func(req *ResponseRequest) {
		req.Instructions = instructions
	}
}

// WithTools sets the tools for the chat completion request.
func WithTools(tools ...ResponseTool) RequestOpt {
	return func(req *ResponseRequest) {
		req.Tools = tools
	}
}
