package jfather

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var input = `abcdefghijklmnopqrstuvwxyz`

func Test_Peeker(t *testing.T) {
	peeker := NewPeekReader(strings.NewReader(input))

	var b byte
	var err error

	for i := 0; i < 30; i++ {
		b, err = peeker.Peek()
		require.NoError(t, err)
		assert.Equal(t, byte('a'), b)
	}

	b, err = peeker.Next()
	require.NoError(t, err)
	assert.Equal(t, byte('a'), b)

	b, err = peeker.Next()
	require.NoError(t, err)
	assert.Equal(t, byte('b'), b)

	b, err = peeker.Peek()
	require.NoError(t, err)
	assert.Equal(t, byte('c'), b)

	data := make([]byte, 5)
	n, err := peeker.Read(data)
	require.NoError(t, err)
	require.Equal(t, 5, n)
	assert.Equal(t, "cdefg", string(data))

	b, err = peeker.Peek()
	require.NoError(t, err)
	assert.Equal(t, byte('h'), b)

	b, err = peeker.Next()
	require.NoError(t, err)
	assert.Equal(t, byte('h'), b)

	data = make([]byte, 20)
	n, err = peeker.Read(data)
	require.NoError(t, err)
	require.Equal(t, 18, n)
	assert.Equal(t, "ijklmnopqrstuvwxyz", string(data[:18]))

	_, err = peeker.Peek()
	require.Error(t, err)

	_, err = peeker.Next()
	require.Error(t, err)

}
