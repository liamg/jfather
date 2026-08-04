package jfather

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// deepArrays returns a valid JSON document of n nested arrays: [[[...]]].
func deepArrays(n int) []byte {
	return []byte(strings.Repeat("[", n) + strings.Repeat("]", n))
}

// deepObjects returns a valid JSON document of n nested objects:
// {"a":{"a":...{}...}}.
func deepObjects(n int) []byte {
	var sb strings.Builder
	for i := 0; i < n-1; i++ {
		sb.WriteString(`{"a":`)
	}
	sb.WriteString("{}")
	sb.WriteString(strings.Repeat("}", n-1))
	return []byte(sb.String())
}

func Test_MaxDepth_DefaultRejectsPathologicalInput(t *testing.T) {
	// The original crash repro: megabytes of unclosed '[' used to overflow
	// the stack and kill the process with an uncatchable fatal error. It
	// must now fail fast with a regular error.
	example := []byte(strings.Repeat("[", 3_000_000))

	var target any
	err := Unmarshal(example, &target)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maximum nesting depth")
}

func Test_MaxDepth_DefaultAllowsDeepButBoundedInput(t *testing.T) {
	var target any
	require.NoError(t, Unmarshal(deepArrays(DefaultMaxDepth), &target))

	err := Unmarshal(deepArrays(DefaultMaxDepth+1), &target)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maximum nesting depth")
}

func Test_MaxDepth_OptionOverridesDefault(t *testing.T) {
	var target any
	require.NoError(t, Unmarshal(deepArrays(5), &target, MaxDepth(5)))
	require.Error(t, Unmarshal(deepArrays(6), &target, MaxDepth(5)))
}

func Test_MaxDepth_CountsObjectsToo(t *testing.T) {
	var target any
	require.NoError(t, Unmarshal(deepObjects(5), &target, MaxDepth(5)))
	require.Error(t, Unmarshal(deepObjects(6), &target, MaxDepth(5)))
}

func Test_MaxDepth_MixedNesting(t *testing.T) {
	example := []byte(`{"a":[{"b":[1]}]}`)

	var target any
	require.NoError(t, Unmarshal(example, &target, MaxDepth(4)))
	require.Error(t, Unmarshal(example, &target, MaxDepth(3)))
}

func Test_MaxDepth_SiblingsDoNotAccumulate(t *testing.T) {
	// Depth tracks nesting, not the total number of containers, so many
	// shallow siblings must parse fine under a tight limit.
	example := []byte(`[[1],[2],[3],[4],{"a":5}]`)

	var target any
	require.NoError(t, Unmarshal(example, &target, MaxDepth(2)))
}

func Test_MaxDepth_InvalidValuesKeepDefault(t *testing.T) {
	var target any
	require.NoError(t, Unmarshal(deepArrays(10), &target, MaxDepth(0)))
	require.NoError(t, Unmarshal(deepArrays(10), &target, MaxDepth(-1)))
}
