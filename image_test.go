package prompts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImage(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Image
	}{
		{
			name: "empty data",
			data: "",
			want: "",
		},
		{
			name: "hello world",
			data: "hello world",
			want: "hello world",
		},
		{
			name: "URL",
			data: "https://example.com/image.png",
			want: "https://example.com/image.png",
		},
		{
			name: "file path",
			data: "/path/to/image.png",
			want: "/path/to/image.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := Image(tt.data)
			require.Equal(t, tt.want, img)
		})
	}
}

func TestImageEncode(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Image
	}{
		{
			name: "empty data",
			data: []byte{},
			want: "",
		},
		{
			name: "hello world",
			data: []byte("hello world"),
			want: "aGVsbG8gd29ybGQ=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var img Image
			img.Encode(tt.data)
			require.Equal(t, tt.want, img)
		})
	}
}

func TestNewImage(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Image
	}{
		{
			name: "empty data",
			data: []byte{},
			want: "",
		},
		{
			name: "hello world",
			data: []byte("hello world"),
			want: "aGVsbG8gd29ybGQ=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := NewImage(tt.data)
			require.Equal(t, tt.want, img)
		})
	}
}
