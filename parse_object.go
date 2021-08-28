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
	// we've hit the end of the object
	if p.swallowIfEqual('}') {
		n.end = p.position
		return n, nil
	}

	// for each key/val
	for {
		key, err := p.parseString()
		if err != nil {
			return nil, err
		}
		n.content = append(n.content, key)

		if err := p.parseWhitespace(); err != nil {
			return nil, err
		}

		if !p.swallowIfEqual(':') {
			return nil, p.makeError("invalid character, expecting ':'")
		}

		valParser := newParser(p.peeker, p.position)
		val, err := valParser.parse()
		if err != nil {
			return nil, err
		}
		p.position = valParser.position
		p.size += valParser.size
		n.content = append(n.content, val)

		// we've hit the end of the object
		if p.swallowIfEqual('}') {
			n.end = p.position
			return n, nil
		}

		if !p.swallowIfEqual(',') {
			return nil, p.makeError("unexpected character - expecting , or }")
		}
	}

}
