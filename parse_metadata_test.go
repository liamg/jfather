package jfather

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParent struct {
	Child testChild `json:"child"`
}

type testChild struct {
	Value       any
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

func (t *testChild) UnmarshalJSONWithMetadata(node Node) error {
	t.StartLine = node.Range().Start.Line
	t.StartColumn = node.Range().Start.Column
	t.EndLine = node.Range().End.Line
	t.EndColumn = node.Range().End.Column
	return node.Decode(&t.Value)
}

func Test_ParseMetadata(t *testing.T) {

	tests := []struct {
		name  string
		input []byte
		want  testParent
	}{
		{
			name: "string",
			input: []byte(`{
 "child": "secret"
}`),
			want: testParent{
				Child: testChild{
					Value:       "secret",
					StartLine:   2,
					StartColumn: 11,
					EndLine:     2,
					EndColumn:   18,
				},
			},
		},
		{
			name: "number",
			input: []byte(`{
 "child": 123
}`),
			want: testParent{
				Child: testChild{
					Value:       int64(123),
					StartLine:   2,
					StartColumn: 11,
					EndLine:     2,
					EndColumn:   13,
				},
			},
		},
		{
			name: "list",
			input: []byte(`{
 "child": [1, 2, 3]
}`),
			want: testParent{
				Child: testChild{
					Value:       []any{int64(1), int64(2), int64(3)},
					StartLine:   2,
					StartColumn: 11,
					EndLine:     2,
					EndColumn:   19,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var parent testParent
			if err := Unmarshal(tt.input, &parent); err != nil {
				panic(err)
			}

			assert.Equal(t, tt.want, parent)

		})
	}
}
