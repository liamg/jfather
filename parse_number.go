package jfather

import (
	"strconv"
	"strings"
)

func (p *parser) parseNumber() (Node, error) {
	n := p.newNode(KindNumber)

	var builder strings.Builder

	if p.swallowIfEqual('-') {
		builder.WriteRune('-')
	}

	if err := p.parseIntegral(&builder); err != nil {
		return nil, err
	}
	hasFraction, err := p.parseFraction(&builder)
	if err != nil {
		return nil, err
	}
	hasExponent, err := p.parseExponent(&builder)
	if err != nil {
		return nil, err
	}

	n.end = p.position

	if hasFraction || hasExponent {
		f, err := strconv.ParseFloat(builder.String(), 64)
		if err != nil {
			return nil, p.makeError("%s", err)
		}
		n.raw = f
		return n, nil
	}

	i, err := strconv.ParseInt(builder.String(), 10, 64)
	if err != nil {
		return nil, p.makeError("%s", err)
	}
	n.raw = i

	return n, nil
}

func (p *parser) parseIntegral(b *strings.Builder) error {
	r, err := p.next()
	if err != nil {
		return err
	}
	if r == '0' {
		r, _ := p.peeker.Peek()
		if r >= '0' && r <= '9' {
			return p.makeError("invalid number")
		}
		b.WriteRune('0')
		return nil
	}

	if r < '1' || r > '9' {
		return p.makeError("invalid number")
	}
	b.WriteRune(r)

	for {
		r, err := p.next()
		if err != nil {
			return nil
		}
		if r < '0' || r > '9' {
			return p.undo()
		}
		b.WriteRune(r)
	}
}

func (p *parser) parseFraction(b *strings.Builder) (bool, error) {
	r, err := p.next()
	if err != nil {
		return false, nil
	}
	if r != '.' {
		return false, p.undo()
	}

	b.WriteRune('.')

	for {
		r, err := p.next()
		if err != nil {
			break
		}
		if r < '0' || r > '9' {
			if err := p.undo(); err != nil {
				return false, err
			}
			break
		}
		b.WriteRune(r)
	}

	if b.String() == "." {
		return false, p.makeError("invalid number - missing digits after decimal point")
	}

	return true, nil
}

func (p *parser) parseExponent(b *strings.Builder) (bool, error) {
	r, err := p.next()
	if err != nil {
		return false, nil
	}
	if r != 'e' && r != 'E' {
		return false, p.undo()
	}

	b.WriteRune(r)

	r, err = p.next()
	if err != nil {
		return true, nil
	}
	hasDigits := r >= '0' && r <= '9'
	if r != '-' && r != '+' && !hasDigits {
		return true, p.undo()
	}
	b.WriteRune(r)

	for {
		r, err := p.next()
		if err != nil {
			break
		}
		if r < '0' || r > '9' {
			if err := p.undo(); err != nil {
				return false, err
			}
			break
		}
		hasDigits = true
		b.WriteRune(r)
	}

	if !hasDigits {
		return false, p.makeError("invalid number - no digits in exponent")
	}

	return true, nil
}
