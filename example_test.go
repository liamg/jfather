package jfather_test

import (
	"fmt"

	"github.com/liamg/jfather"
)

func ExampleAllowComments() {
	input := []byte(`{
	// the service to deploy
	"name": "example" /* block comment */
}`)

	var target struct {
		Name string `json:"name"`
	}
	if err := jfather.Unmarshal(input, &target, jfather.AllowComments()); err != nil {
		panic(err)
	}

	fmt.Println(target.Name)
	// Output: example
}

func ExampleAllowTrailingCommas() {
	input := []byte(`{ "tags": [ "a", "b", ], }`)

	var target struct {
		Tags []string `json:"tags"`
	}
	if err := jfather.Unmarshal(input, &target, jfather.AllowTrailingCommas()); err != nil {
		panic(err)
	}

	fmt.Println(target.Tags)
	// Output: [a b]
}
