// Package jfather is a JSON parser and decoder that provides metadata about the parsed JSON.
package jfather

import "bytes"

// Unmarshaler is an interface that can be implemented by types that can unmarshal a JSON description of themselves.
type Unmarshaler interface {
	UnmarshalJSONWithMetadata(node Node) error
}

// Unmarshal unmarshals the given JSON data into the target value.
// It will attempt to use UnmarshalJSONWithMetadata if the target implements it, otherwise it will use the default json unmarshaler.
func Unmarshal(data []byte, target any) error {
	node, err := newParser(NewPeekReader(bytes.NewReader(data)), Position{1, 1}).parse()
	if err != nil {
		return err
	}
	return node.Decode(target)
}
