package jfather

// Node represents a node in the AST.
type Node interface {
	Range() Range
	Decode(target any) error
	Kind() Kind
	Content() []Node
}

// Range represents a range of positions in the source.
type Range struct {
	Start Position
	End   Position
}

// Position represents a position in the source code. Note that both lines and columns are 1-indexed.
type Position struct {
	Line   int
	Column int
}

type node struct {
	raw     any
	content []Node
	start   Position
	end     Position
	kind    Kind
}

func (n *node) Range() Range {
	return Range{
		Start: n.start,
		End: Position{
			Column: n.end.Column - 1,
			Line:   n.end.Line,
		},
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
