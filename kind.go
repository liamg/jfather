package jfather

// Kind represents the type of a node.
type Kind uint8

const (
	// KindUnknown represents an unknown kind.
	KindUnknown Kind = iota
	// KindNull represents a null value.
	KindNull
	// KindNumber represents a number value.
	KindNumber
	// KindString represents a string value.
	KindString
	// KindBoolean represents a boolean value.
	KindBoolean
	// KindArray represents an array value.
	KindArray
	// KindObject represents an object value.
	KindObject
)
