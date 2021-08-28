package jfather

import "fmt"

func (n *node) Decode(target interface{}) error {
	switch n.kind {
	case KindObject:
		return n.decodeObject(target)
	case KindArray:
		return n.decodeArray(target)
	case KindString:
		return n.decodeString(target)
	case KindNumber:
		return n.decodeNumber(target)
	case KindBoolean:
		return n.decodeBoolean(target)
	case KindNull:
		return n.decodeNull(target)
	case KindUnknown:
		return fmt.Errorf("cannot decode unknown kind")
	default:
		return fmt.Errorf("decoding of kind 0x%x not yet implemented", n.kind)
	}
}
