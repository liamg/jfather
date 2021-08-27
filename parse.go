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

	if err := p.eatWhitespace(); err != nil {
		return nil, err
	}

	c, err := p.peeker.Peek()
	if err != nil {
		return nil, err
	}

	switch c {
	case '{':
		return p.parseObject()
	default:
		return nil, fmt.Errorf("unexpected character '%c'", c)

	}
}

func (p *parser) next() (byte, error) {
	b, err := p.peeker.Next()
	if err != nil {
		return 0, err
	}
	p.column++
	p.size++
	return b, nil
}

func (p *parser) eatWhitespace() error {
	for {
		b, err := p.peeker.Peek()
		if err != nil {
			return err
		}
		switch b {
		case 0x0d, 0x20, 0x09:
			p.column++
		case 0x0a:
			p.column = 1
			p.line++
		default:
			return nil
		}
		if _, err := p.peeker.Next(); err != nil {
			return err
		}
		p.size++
	}
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

func (p *parser) parseString() (Node, error) {

	if err := p.eatWhitespace(); err != nil {
		return nil, err
	}

	n := p.newNode(KindString)

	b, err := p.next()
	if err != nil {
		return nil, err
	}

	if b != byte('"') {
		return nil, p.makeError("expecting string delimiter")
	}

	n.data = []byte{'"'}

	var inEscape bool

	for {
		c, err := p.next()
		if err != nil {
			return nil, err
		}
		n.data = append(n.data, c)
		if inEscape {
			inEscape = false
		} else {
			switch c {
			case '\\':
				inEscape = true
			case '"':
				n.end.Line = p.line
				n.end.Column = p.column
				return n, nil
			}
		}

	}
}

func (p *parser) swallowIfEqual(b byte) bool {
	c, err := p.peeker.Peek()
	if err != nil {
		return false
	}
	if c != b {
		return false
	}
	_, _ = p.next()
	return true
}

func (p *parser) parseObject() (Node, error) {

	if err := p.eatWhitespace(); err != nil {
		return nil, err
	}

	n := p.newNode(KindObject)
	c, err := p.next()
	if err != nil {
		return nil, err
	}

	if c != '{' {
		return nil, p.makeError("expecting object delimiter")
	}
	n.data = append(n.data, c)
	if err := p.eatWhitespace(); err != nil {
		return nil, err
	}

	// for each key/val
	for {

		// we've hit the end of the object
		if p.swallowIfEqual('}') {
			n.data = append(n.data, c)
			n.end.Column = p.column
			n.end.Line = p.line
			return n, nil
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		if err := p.eatWhitespace(); err != nil {
			return nil, err
		}

		if !p.swallowIfEqual(':') {
			return nil, p.makeError("invalid character, expecting ':'")
		}

		valParser := newParser(p.peeker, p.line, p.column)
		val, err := valParser.parse()
		if err != nil {
			return nil, err
		}
		p.line = valParser.line
		p.column = valParser.column
		p.size += valParser.size

		_, _ = key, val
	}

}
