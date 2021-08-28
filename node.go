package jfather

type Node interface {
	Range() Range
	Decode(target interface{}) error
	Kind() Kind
	Content() []Node
}

type Range struct {
	Start Position
	End   Position
}

type Position struct {
	Line   int
	Column int
}

type node struct {
	raw     interface{}
	start   Position
	end     Position
	kind    Kind
	content []Node
}

func (n *node) Range() Range {
	return Range{
		Start: n.start,
		End:   n.end,
	}
}

func (n *node) End() Position {
	return n.end
}

func (n *node) Kind() Kind {
	return n.kind
}

func (n *node) Content() []Node {
	return n.content
}
