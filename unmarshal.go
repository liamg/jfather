package jfather

import "bytes"

type Unmarshaller interface {
	UnmarshalJFather(node Node) error
}

func Unmarshal(data []byte, target interface{}) error {
	node, err := newParser(NewPeekReader(bytes.NewReader(data)), Position{1, 1}).parse()
	if err != nil {
		return err
	}
	if unmarshaller, ok := target.(Unmarshaller); ok {
		return unmarshaller.UnmarshalJFather(node)
	}
	return node.Decode(target)
}
