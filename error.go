package prompts

import "encoding/json"

// Error represents an error returned by the API.
type Error struct {
	JSON struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.JSON.Message
}

// UnmarshalJSON unmarshals the error from JSON.
func (e *Error) UnmarshalJSON(data []byte) error {
	var err struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	}

	if err := json.Unmarshal(data, &err); err != nil {
		return err
	}

	e.JSON.Code = err.Code
	e.JSON.Message = err.Message
	e.JSON.Type = err.Type

	return nil
}

// FromBody unmarshals the error from JSON.
func FromBody(data []byte) error {
	var err Error
	if err := json.Unmarshal(data, &err); err != nil {
		return err
	}

	return &err
}
