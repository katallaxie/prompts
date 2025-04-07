package prompts

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
