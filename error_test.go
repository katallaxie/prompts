package prompts_test

import (
	"testing"

	"github.com/katallaxie/prompts"
	"github.com/stretchr/testify/require"
)

func TestPromptErro(t *testing.T) {
	tests := []struct {
		name string
		data string
		want prompts.PromptError
	}{
		{
			name: "empty data",
			data: "",
			want: prompts.PromptError{},
		},
		{
			name: "invalid JSON",
			data: "{invalid json}",
			want: prompts.PromptError{},
		},
		{
			name: "valid error",
			data: `{"error":{"code":400,"message":"Bad Request","type":"invalid_request_error"}}`,
			want: prompts.PromptError{
				JSON: struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Type    string `json:"type"`
				}{
					Code:    400,
					Message: "Bad Request",
					Type:    "invalid_request_error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err prompts.PromptError
			err.UnmarshalJSON([]byte(tt.data))
			require.Equal(t, tt.want, err)
		})
	}
}
