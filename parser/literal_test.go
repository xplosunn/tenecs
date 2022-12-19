package parser_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"testing"
)

func TestLiteralFoldNil(t *testing.T) {
	assert.Equal(t, "", parser.LiteralFold[string](nil, nil, nil, nil, nil))
	assert.Equal(t, "", parser.LiteralFold(
		nil,
		func(arg float64) string { panic("unexpected call") },
		func(arg int) string { panic("unexpected call") },
		func(arg string) string { panic("unexpected call") },
		func(arg bool) string { panic("unexpected call") },
	))
}
