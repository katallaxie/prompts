package prompts

// ToolFunction represents a function to call when a tool is used.
type ToolFunction struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is the description of the function.
	Description string `json:"description"`
	// Parameters is the parameters to pass to the function.
	Parameters []ToolParameters `json:"parameters"`
}

// ToolProperty represents a property for a tool function.
type ToolProperty struct {
	// Type is the type of the property.
	Type string `json:"type"`
	// Description is the description of the property.
	Description string `json:"description"`
}

// ToolParameter represents a parameter for a tool function.
type ToolParameters struct {
	// Type is the type of the parameter.
	Type string `json:"type"`
	// Description is the description of the parameter.
	Description string `json:"description"`
	// Properties is the properties of the parameter.
	Properties map[string]ToolProperty `json:"properties"`
	// Required indicates whether the parameter is required.
	Required []string `json:"required"`
}

// ToolCall represents a tool call that can be used in a chat completion response.
type ToolCall struct {
	// Function is the function to call when the tool is used.
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents a tool call function that can be used in a chat completion response.
type ToolCallFunction struct {
	// Name is the name of the function to call when the tool is used.
	Name string `json:"name"`
	// Arguments is the arguments to pass to the function when the tool is used.
	Arguments map[string]interface{} `json:"arguments"`
}

// Tool represents a tool that can be used in a chat completion request.
type Tool struct {
	// Type is the type of the tool.
	Type string `json:"type"`
	// Function is the function to call when the tool is used.
	Function ToolFunction `json:"function"`
}
