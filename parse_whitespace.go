package jfather

import "io"

func (p *parser) parseWhitespace() error {
	for {
		b, err := p.peeker.Peek()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch b {
		case 0x0d, 0x20, 0x09, 0x0a:
			if _, err := p.next(); err != nil {
				return err
			}
			if b == 0x0a {
				p.position.Column = 1
				p.position.Line++
			}
		case '/':
			if !p.allowComments {
				return nil
			}
			if err := p.skipComment(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}
