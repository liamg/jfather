package jfather

import (
	"fmt"
	"io"
)

type PeekReader struct {
	peeked     bool
	buffer     byte
	underlying io.Reader
}

func NewPeekReader(reader io.Reader) *PeekReader {
	return &PeekReader{
		underlying: reader,
	}
}

func (r *PeekReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("cannot read into zero-length slice")
	}
	var offset int
	if r.peeked {
		b[0] = r.buffer
		r.peeked = false
		offset = 1
	}
	nn, err := r.underlying.Read(b[offset:])
	return offset + nn, err
}

func (r *PeekReader) Next() (byte, error) {
	in := make([]byte, 1)
	n, err := r.Read(in)
	if n > 0 {
		return in[0], nil
	}
	return 0, err
}

func (r *PeekReader) Peek() (byte, error) {
	b, err := r.Next()
	if err != nil {
		return 0, err
	}
	r.peeked = true
	r.buffer = b
	return b, nil
}
