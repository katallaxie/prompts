package prompts

import (
	"encoding/base64"
)

// Image is a type that represents an image.
type Image string

// Encode encodes the image into a string.
func (i *Image) Encode(data []byte) string {
	*i = Image(base64.StdEncoding.EncodeToString(data))
	return string(*i)
}

// String returns the string representation of the image.
func (i Image) String() string {
	return string(i)
}

// NewImage creates a new image from the given data.
func NewImage(data []byte) Image {
	var img Image
	img.Encode(data)
	return img
}
