package jfather

import (
	"fmt"
)

type parser struct {
	line   int
	column int
	size   int
	peeker *PeekReader
}

func newParser(p *PeekReader, line int, column int) *parser {
	return &parser{
		line:   line,
		column: column,
		peeker: p,
	}
}

func (p *parser) parse() (Node, error) {

	if err := p.parseWhitespace(); err != nil {
		return nil, err
	}

	c, err := p.peeker.Peek()
	if err != nil {
		return nil, err
	}

	switch c {
	case '"':
		return p.parseString()
	case '{':
		return p.parseObject()
	default:
		return nil, fmt.Errorf("unexpected character '%c'", c)

	}
}

func (p *parser) next() (rune, error) {
	b, err := p.peeker.Next()
	if err != nil {
		return 0, err
	}
	p.column++
	p.size++
	return b, nil
}

func (p *parser) makeError(format string, args ...interface{}) error {
	return fmt.Errorf(
		"Error at line %d, column %d: %s",
		p.line,
		p.column,
		fmt.Sprintf(format, args...),
	)
}

func (p *parser) newNode(k Kind) *node {
	return &node{
		start: Range{
			Line:   p.line,
			Column: p.column,
		},
		kind: k,
	}
}

func (p *parser) swallowIfEqual(r rune) bool {
	c, err := p.peeker.Peek()
	if err != nil {
		return false
	}
	if c != r {
		return false
	}
	_, _ = p.next()
	return true
}
