package parser_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestParseString(t *testing.T) {
	testCases := testcode.GetAll()
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := parser.ParseString(testCase.Content)
			assert.NoError(t, err)
		})
	}
}
