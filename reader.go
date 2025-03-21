package prompts

import (
	"io"
)

// CompletionReader is a completion reader.
type CompletionReader interface {
	io.Reader
}

// CompletionReaderrImpl is a completion reader implementation.
type CompletionReaderImpl struct {
	s        []byte
	i        int64
	prevRune int
}

// NewCompletionReader returns a new completion reader.
func NewCompletionReader(completion ...CompletionChoice) CompletionReader {
	r := new(CompletionReaderImpl)

	for _, c := range completion {
		r.s = append(r.s, []byte(c.Message.GetContent())...)
	}

	return &CompletionReaderImpl{s: r.s}
}

// Read reads a message from the stream.
func (r *CompletionReaderImpl) Read(b []byte) (int, error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n := copy(b, r.s[r.i:])
	r.i += int64(n)

	return n, nil
}
