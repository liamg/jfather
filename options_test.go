package jfather

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type optTarget struct {
	Name    string   `json:"name"`
	Balance float64  `json:"balance"`
	Tags    []string `json:"tags"`
}

func Test_AllowComments_LineComments(t *testing.T) {
	example := []byte(`{
	// leading line comment
	"name": "testing", // trailing line comment
	"balance": 3.14
	// comment before close
}`)

	var target optTarget
	require.NoError(t, Unmarshal(example, &target, AllowComments()))
	assert.Equal(t, "testing", target.Name)
	assert.Equal(t, 3.14, target.Balance)
}

func Test_AllowComments_BlockComments(t *testing.T) {
	example := []byte(`{
	/* block before key */ "name": "testing",
	"balance": /* inline */ 3.14,
	/* multi
	   line
	   block */
	"tags": [ "a", /* mid-array */ "b" ]
}`)

	var target optTarget
	require.NoError(t, Unmarshal(example, &target, AllowComments()))
	assert.Equal(t, "testing", target.Name)
	assert.Equal(t, 3.14, target.Balance)
	assert.Equal(t, []string{"a", "b"}, target.Tags)
}

func Test_AllowComments_DoesNotTouchStringContents(t *testing.T) {
	// // and /* inside string values must be preserved verbatim.
	example := []byte(`{ "name": "http://example.com/*not a comment*/" }`)

	var target optTarget
	require.NoError(t, Unmarshal(example, &target, AllowComments()))
	assert.Equal(t, "http://example.com/*not a comment*/", target.Name)
}

func Test_Comments_RejectedByDefault(t *testing.T) {
	example := []byte(`{ "name": "testing" // nope
}`)
	var target optTarget
	assert.Error(t, Unmarshal(example, &target), "strict mode must reject comments")
}

func Test_AllowComments_UnterminatedBlock(t *testing.T) {
	example := []byte(`{ "name": "testing" /* never closed `)
	var target optTarget
	err := Unmarshal(example, &target, AllowComments())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unterminated block comment")
}

func Test_AllowComments_LoneSlashIsError(t *testing.T) {
	example := []byte(`{ "name": / "testing" }`)
	var target optTarget
	assert.Error(t, Unmarshal(example, &target, AllowComments()))
}

func Test_AllowTrailingCommas_Object(t *testing.T) {
	example := []byte(`{
	"name": "testing",
	"balance": 3.14,
}`)
	var target optTarget
	require.NoError(t, Unmarshal(example, &target, AllowTrailingCommas()))
	assert.Equal(t, "testing", target.Name)
	assert.Equal(t, 3.14, target.Balance)
}

func Test_AllowTrailingCommas_Array(t *testing.T) {
	example := []byte(`{ "tags": [ "a", "b", ] }`)
	var target optTarget
	require.NoError(t, Unmarshal(example, &target, AllowTrailingCommas()))
	assert.Equal(t, []string{"a", "b"}, target.Tags)
}

func Test_TrailingCommas_RejectedByDefault(t *testing.T) {
	example := []byte(`{ "tags": [ "a", "b", ] }`)
	var target optTarget
	assert.Error(t, Unmarshal(example, &target), "strict mode must reject trailing commas")
}

func Test_Options_Combined(t *testing.T) {
	// Comments and trailing commas together, mimicking a real ARM-style
	// document.
	example := []byte(`{
	// resource config
	"name": "testing",
	"tags": [
		"a", // first
		"b", /* last */
	],
	"balance": 3.14,
}`)
	var target optTarget
	require.NoError(t, Unmarshal(example, &target, AllowComments(), AllowTrailingCommas()))
	assert.Equal(t, "testing", target.Name)
	assert.Equal(t, []string{"a", "b"}, target.Tags)
	assert.Equal(t, 3.14, target.Balance)
}

// Metadata (line/column) must stay correct after a multi-line block
// comment shifts subsequent tokens down.
func Test_AllowComments_PreservesMetadataAfterBlock(t *testing.T) {
	example := []byte("{\n/* a\nb */\n\"name\": \"testing\"\n}")

	var meta metaObject
	require.NoError(t, Unmarshal(example, &meta, AllowComments()))
	// "name" sits on line 4 after the 3-line block comment.
	assert.Equal(t, 4, meta.nameLine, "line number must account for newlines inside the block comment")
}

type metaObject struct {
	nameLine int
}

func (m *metaObject) UnmarshalJSONWithMetadata(node Node) error {
	content := node.Content()
	for i := 0; i+1 < len(content); i += 2 {
		var key string
		if err := content[i].Decode(&key); err != nil {
			return err
		}
		if key == "name" {
			m.nameLine = content[i+1].Range().Start.Line
		}
	}
	return nil
}

func Test_AllowUnescapedControlChars_Newline(t *testing.T) {
	// A raw newline inside a string value — as ARM templates do when
	// breaking a long expression across multiple lines.
	input := []byte("{ \"expr\": \"[concat('a',\n'b')]\" }")
	var target struct {
		Expr string `json:"expr"`
	}
	require.NoError(t, Unmarshal(input, &target, AllowUnescapedControlChars()))
	assert.Equal(t, "[concat('a',\n'b')]", target.Expr)
}

func Test_AllowUnescapedControlChars_Tab(t *testing.T) {
	input := []byte("{ \"v\": \"a\tb\" }")
	var target struct {
		V string `json:"v"`
	}
	require.NoError(t, Unmarshal(input, &target, AllowUnescapedControlChars()))
	assert.Equal(t, "a\tb", target.V)
}

func Test_UnescapedControlChars_RejectedByDefault(t *testing.T) {
	input := []byte("{ \"v\": \"a\nb\" }")
	var target struct {
		V string `json:"v"`
	}
	err := Unmarshal(input, &target)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid unescaped character")
}

// A literal newline inside a string must advance the line counter so
// metadata for subsequent tokens stays accurate.
func Test_AllowUnescapedControlChars_TracksLineNumbers(t *testing.T) {
	input := []byte("{\n\"a\": \"x\ny\",\n\"name\": \"here\"\n}")
	var meta metaObject
	require.NoError(t, Unmarshal(input, &meta, AllowUnescapedControlChars()))
	// "a" is on line 2; its value spans a raw newline; "name" is on line 4.
	assert.Equal(t, 4, meta.nameLine)
}
