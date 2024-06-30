package codegen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestCode(t *testing.T) {
	for _, testCode := range testcode.GetAll() {
		t.Run(testCode.Name, func(t *testing.T) {
			parsed, err := parser.ParseString(testCode.Content)
			assert.NoError(t, err)

			typed, err := typer.TypecheckSingleFile(*parsed)
			assert.NoError(t, err)

			generated := codegen.GenerateProgramTest(typed)

			golang.RunCodeUnlessCached(t, generated)
		})
	}
}
