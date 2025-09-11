package jfather

import (
	"bufio"
	"io"
)

// PeekReader is a reader that allows peeking at the next rune.
type PeekReader struct {
	underlying *bufio.Reader
}

// NewPeekReader returns a new PeekReader.
func NewPeekReader(reader io.Reader) *PeekReader {
	return &PeekReader{
		underlying: bufio.NewReader(reader),
	}
}

// Next returns the next rune.
func (r *PeekReader) Next() (rune, error) {
	c, _, err := r.underlying.ReadRune()
	return c, err
}

// Undo undoes the last Next call.
func (r *PeekReader) Undo() error {
	return r.underlying.UnreadRune()
}

// Peek returns the next rune without advancing the reader.
func (r *PeekReader) Peek() (rune, error) {
	c, _, err := r.underlying.ReadRune()
	if err != nil {
		return 0, err
	}
	if err := r.underlying.UnreadRune(); err != nil {
		return 0, err
	}
	return c, nil
}
