package perplexity

import (
	"context"

	"github.com/katallaxie/prompts"
	"github.com/katallaxie/prompts/openai"
)

type (
	ResponseRequest            = openai.ResponseRequest
	Response                   = openai.Response
	ResponseInput              = openai.ResponseInput
	ResponseTool               = openai.ResponseTool
	ResponseMessageContent     = openai.ResponseMessageContent
	ResponseMessageContentText = openai.ResponseMessageContentText
	ResponseFunctionTool       = openai.ResponseFunctionTool
	ResponseFunctionDefinition = openai.ResponseFunctionDefinition
	ResponseFunctionParameters = openai.ResponseFunctionParameters
	ResponseFunctionProperties = openai.ResponseFunctionProperties
	Role                       = openai.Role
)

// Role constants for the Perplexity API.
const (
	RoleAgent     Role = "agent"
	RoleNone      Role = "none"
	RoleSystem    Role = "system"
	RoleTool      Role = "tool"
	RoleUser      Role = "user"
	RoleDeveloper Role = "developer"
)

// Ollama is a struct that implements the Prompter interface for the Ollama API.
type Ollama[I *ResponseRequest, O *Response] struct {
	client *prompts.Client
}

// New creates a new Ollama with the given client.
func New(client *prompts.Client) prompts.Prompter[*ResponseRequest, *Response] {
	base := client.New().Base(DefaultURL)

	return &Ollama[*ResponseRequest, *Response]{client: base}
}

// Respond sends a chat completion request and returns the response.
func (p *Ollama[I, O]) Respond(ctx context.Context, req I) (O, error) {
	res := &Response{}

	_, err := p.client.New().Post("responses").BodyJSON(req).ReceiveSuccess(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DefaultURL is the default endpoint for the Ollama API.
const DefaultURL = "http://localhost:11434/v1/"

// DefaultModel is the default model for the Ollama API.
const DefaultModel = "qwen3:8b"
