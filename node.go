package main

type Node interface {
	Offset() int
	Start() Range
	End() Range
	Decode(target interface{}) error
}

type Range struct {
	Line   int
	Column int
}

type node struct {
	data   []byte
	start  Range
	end    Range
	offset int
}

func (n *node) Decode(target interface{}) error {
	return decode(n, target)
}

func (n *node) Offset() int {
	return n.offset
}
func (n *node) Start() Range {
	return n.start
}
func (n *node) End() Range {
	return n.end
}
