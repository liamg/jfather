package jfather

func (p *parser) parseWhitespace() error {
	for {
		b, err := p.peeker.Peek()
		if err != nil {
			return err
		}
		switch b {
		case 0x0d, 0x20, 0x09:
			_, err := p.next()
			return err
		case 0x0a:
			p.column = 1
			p.line++
		default:
			return nil
		}
	}
}
