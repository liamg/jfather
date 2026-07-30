// Package jfather is a JSON parser and decoder that provides metadata about the parsed JSON.
package jfather

import "bytes"

// Unmarshaler is an interface that can be implemented by types that can unmarshal a JSON description of themselves.
type Unmarshaler interface {
	UnmarshalJSONWithMetadata(node Node) error
}

// Unmarshal unmarshals the given JSON data into the target value.
// It will attempt to use UnmarshalJSONWithMetadata if the target implements it, otherwise it will use the default json unmarshaler.
//
// By default the input must be strict JSON. Pass options such as
// AllowComments or AllowTrailingCommas to accept common JSON dialects.
func Unmarshal(data []byte, target any, opts ...Option) error {
	node, err := newParser(NewPeekReader(bytes.NewReader(data)), Position{1, 1}, opts...).parse()
	if err != nil {
		return err
	}
	return node.Decode(target)
}
