package jfather

type Node interface {
	Start() Range
	End() Range
	Decode(target interface{}) error
	Kind() Kind
	Content() []Node
}

type Range struct {
	Line   int
	Column int
}

type node struct {
	raw     interface{}
	start   Range
	end     Range
	kind    Kind
	content []Node
}

func (n *node) Decode(target interface{}) error {
	return decode(n, target)
}

func (n *node) Start() Range {
	return n.start
}

func (n *node) End() Range {
	return n.end
}

func (n *node) Kind() Kind {
	return n.kind
}

func (n *node) Content() []Node {
	return n.content
}
