package prompts

import "encoding/json"

// PromptError represents an error returned by the API.
type PromptError struct {
	JSON struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Error returns the error message.
func (e *PromptError) Error() string {
	return e.JSON.Message
}

// UnmarshalJSON unmarshals the error from JSON.
func (e *PromptError) UnmarshalJSON(data []byte) error {
	var err struct {
		JSON struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}

	if err := json.Unmarshal(data, &err); err != nil {
		return err
	}

	e.JSON.Code = err.JSON.Code
	e.JSON.Message = err.JSON.Message
	e.JSON.Type = err.JSON.Type

	return nil
}
