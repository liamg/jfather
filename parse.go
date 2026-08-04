package jfather

import (
	"fmt"
	"io"
)

// DefaultMaxDepth is the maximum nesting depth of arrays and objects
// accepted by Unmarshal unless overridden with the MaxDepth option. The
// parser (and any code that later walks the resulting tree) recurses once
// per nesting level, so depth must be bounded to prevent maliciously deep
// input from overflowing the stack. The value matches encoding/json's
// nesting limit.
const DefaultMaxDepth = 10000

type parser struct {
	peeker                 *PeekReader
	position               Position
	size                   int
	depth                  int
	maxDepth               int
	allowComments          bool
	allowTrailingCommas    bool
	allowUnescapedControls bool
}

func newParser(p *PeekReader, pos Position, opts ...Option) *parser {
	parser := &parser{
		position: pos,
		peeker:   p,
		maxDepth: DefaultMaxDepth,
	}
	for _, opt := range opts {
		opt(parser)
	}
	return parser
}

func (p *parser) parse() (Node, error) {
	return p.parseElement()
}

func (p *parser) parseElement() (Node, error) {
	if err := p.parseWhitespace(); err != nil {
		return nil, err
	}
	n, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	if err := p.parseWhitespace(); err != nil {
		return nil, err
	}
	return n, nil
}

func (p *parser) parseValue() (Node, error) {
	c, err := p.peeker.Peek()
	if err != nil {
		return nil, err
	}

	switch c {
	case '"':
		return p.parseString()
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case 'n':
		return p.parseNull()
	case 't', 'f':
		return p.parseBoolean()
	default:
		if c == '-' || (c >= '0' && c <= '9') {
			return p.parseNumber()
		}
		return nil, fmt.Errorf("unexpected character '%c'", c)
	}
}

func (p *parser) next() (rune, error) {
	b, err := p.peeker.Next()
	if err != nil {
		return 0, err
	}
	p.position.Column++
	p.size++
	return b, nil
}

func (p *parser) undo() error {
	if err := p.peeker.Undo(); err != nil {
		return err
	}
	p.position.Column--
	p.size--
	return nil
}

func (p *parser) makeError(format string, args ...any) error {
	return fmt.Errorf(
		"error at line %d, column %d: %s",
		p.position.Line,
		p.position.Column,
		fmt.Sprintf(format, args...),
	)
}

// enterNesting is called on entry to parseObject/parseArray, the two
// productions that recurse. It must run before any recursive descent so
// over-deep input is rejected while the stack is still shallow. Callers
// must pair it with exitNesting so depth reflects nesting, not total
// container count.
func (p *parser) enterNesting() error {
	p.depth++
	if p.depth > p.maxDepth {
		return p.makeError("maximum nesting depth of %d exceeded", p.maxDepth)
	}
	return nil
}

func (p *parser) exitNesting() {
	p.depth--
}

func (p *parser) newNode(k Kind) *node {
	return &node{
		start: p.position,
		kind:  k,
	}
}

// skipComment consumes a // line comment or /* block comment. It is
// called from parseWhitespace when allowComments is set and the next rune
// is '/'. A line comment stops at (but does not consume) the terminating
// newline, leaving parseWhitespace to account for it. A '/' not followed
// by '/' or '*' is malformed — comments are the only place '/' is legal
// outside a string.
func (p *parser) skipComment() error {
	// consume the leading '/'
	if _, err := p.next(); err != nil {
		return err
	}
	c, err := p.next()
	if err != nil {
		return err
	}

	switch c {
	case '/':
		for {
			b, err := p.peeker.Peek()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if b == 0x0a {
				return nil
			}
			if _, err := p.next(); err != nil {
				return err
			}
		}
	case '*':
		for {
			b, err := p.next()
			if err != nil {
				if err == io.EOF {
					return p.makeError("unterminated block comment")
				}
				return err
			}
			if b == 0x0a {
				p.position.Column = 1
				p.position.Line++
				continue
			}
			if b == '*' {
				if p.swallowIfEqual('/') {
					return nil
				}
			}
		}
	default:
		return p.makeError("unexpected character '%c' after '/'", c)
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
