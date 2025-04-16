package codegen_golang_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

// TODO FIXME figure out what tests are running here and if that makes sense
func TestCode(t *testing.T) {
	for _, testCode := range testcode.GetAll() {
		t.Run(testCode.Name, func(t *testing.T) {
			parsed, err := parser.ParseString(testCode.Content)
			assert.NoError(t, err)

			typed, err := typer.TypecheckSingleFile(*parsed)
			assert.NoError(t, err)

			generated := codegen_golang.GenerateProgramTest(typed, codegen.FindTests(typed))

			golang.RunCodeUnlessCached(t, generated)
		})
	}
}
