package jfather

func (p *parser) parseObject() (Node, error) {

	if err := p.parseWhitespace(); err != nil {
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
	if err := p.parseWhitespace(); err != nil {
		return nil, err
	}

	// for each key/val
	for {

		// we've hit the end of the object
		if p.swallowIfEqual('}') {
			n.end.Column = p.column
			n.end.Line = p.line
			return n, nil
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		if err := p.parseWhitespace(); err != nil {
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
