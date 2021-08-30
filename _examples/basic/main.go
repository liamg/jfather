package main

import (
	"fmt"

	"github.com/liamg/jfather"
)

type ExampleParent struct {
	Child *ExampleChild `json:"child"`
}

type ExampleChild struct {
	Name   string
	Line   int
	Column int
}

func (t *ExampleChild) UnmarshalJSONWithMetadata(node jfather.Node) error {
	t.Line = node.Range().Start.Line
	t.Column = node.Range().Start.Column
	return node.Decode(&t.Name)
}

func main() {
	input := []byte(`{
	"child": "secret"
}`)
	var parent ExampleParent
	if err := jfather.Unmarshal(input, &parent); err != nil {
		panic(err)
	}

	fmt.Printf("Child value is at line %d, column %d, and is set to '%s'\n",
		parent.Child.Line, parent.Child.Column, parent.Child.Name)

	// outputs:
	//  Child value is at line 2, column 12, and is set to 'secret'
}
