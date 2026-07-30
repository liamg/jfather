package jfather

func (p *parser) parseArray() (Node, error) {
	n := p.newNode(KindArray)

	c, err := p.next()
	if err != nil {
		return nil, err
	}

	if c != '[' {
		return nil, p.makeError("expecting object delimiter")
	}
	if err := p.parseWhitespace(); err != nil {
		return nil, err
	}
	// we've hit the end of the object
	if p.swallowIfEqual(']') {
		n.end = p.position
		return n, nil
	}

	n.content = make([]Node, 0, 8)

	// for each element
	for {

		if err := p.parseWhitespace(); err != nil {
			return nil, err
		}

		val, err := p.parseElement()
		if err != nil {
			return nil, err
		}
		n.content = append(n.content, val)

		// we've hit the end of the array
		if p.swallowIfEqual(']') {
			n.end = p.position
			return n, nil
		}

		if !p.swallowIfEqual(',') {
			return nil, p.makeError("unexpected character - expecting , or ]")
		}

		// A trailing comma is allowed to precede the closing ']' when the
		// option is set.
		if p.allowTrailingCommas {
			if err := p.parseWhitespace(); err != nil {
				return nil, err
			}
			if p.swallowIfEqual(']') {
				n.end = p.position
				return n, nil
			}
		}
	}
}
