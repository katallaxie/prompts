package prompts

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
	MaxTokens int `json:"max_tokens"`
	// Temperature is the sampling temperature
	Temperature *float32 `json:"temperature"`
	// Stream is a flag to enable streaming
	Stream *bool `json:"stream,omitempty"`
	// TopP is the nucleus sampling parameter
	TopP *float64 `json:"top_p,omitempty"`
	// TopK is the number of top tokens to sample from
	TopK *int `json:"top_k,omitempty"`
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
