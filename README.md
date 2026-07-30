# jfather

Parse JSON with line numbers and more!

This is a JSON parsing module that provides additional information during the unmarshalling process, such as line numbers, columns etc.

You can use jfather to unmarshal JSON just like the `encoding/json` package, and add your own unmarshalling functionality to gather metadata by implementing the `jfather.Unmarshaler` interface. This requires a single method with the signature `UnmarshalJSONWithMetadata(node jfather.Node) error`. A full example is below.

You should not use this package unless you need the line/column metadata, as unmarshalling is typically much slower than the `encoding/json` package:

```
BenchmarkUnmarshal_JFather-11    	  120945	      9401 ns/op	   14016 B/op	     176 allocs/op
BenchmarkUnmarshal_Traditional-11     326814	      3699 ns/op	    2552 B/op	      56 allocs/op
```

## Full Example

```golang
package main

import (
	"fmt"

	"github.com/liamg/jfather"
)

type ExampleParent struct {
	Child ExampleChild `json:"child"`
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
	//  Child value is at line 2, column 11, and is set to 'secret'
}
```

## Options

By default `Unmarshal` accepts strict [RFC 8259](https://datatracker.ietf.org/doc/html/rfc8259) JSON. You can pass options to accept common non-standard JSON dialects. Both options are opt-in, so the default behaviour is unchanged.

- `AllowComments()` — allow `//` line comments and `/* ... */` block comments anywhere whitespace is permitted.
- `AllowTrailingCommas()` — allow a single trailing comma before a closing `]` or `}`.
- `AllowUnescapedControlChars()` — allow literal control characters (raw newlines, tabs, etc.) inside string values, instead of requiring them to be escaped. Mirrors Python's `json.loads(..., strict=False)`.

These are handy for formats such as [JSONC](https://code.visualstudio.com/docs/languages/json#_json-with-comments), `tsconfig.json`, and Azure ARM/Bicep deployment templates, all of which permit comments (and, in ARM's case, long expressions broken across multiple lines inside a single string).

```golang
input := []byte(`{
	// the service to deploy
	"name": "example",
	"tags": [ "a", "b", ], /* trailing comma */
}`)

var target map[string]any
if err := jfather.Unmarshal(input, &target,
	jfather.AllowComments(),
	jfather.AllowTrailingCommas(),
); err != nil {
	panic(err)
}
```
