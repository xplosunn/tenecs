package codegen_golang_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestErrorWithLineNumber(t *testing.T) {
	program := `package test

import tenecs.int.plus
import tenecs.test.UnitTest
import tenecs.test.UnitTestKit

_ := UnitTest("plus", (testkit: UnitTestKit): Void => {
  testkit.assert.equal (3, plus(1, 1))
})
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramTest(typed, codegen.FindTests(typed))

	result := golang.RunCodeUnlessCached(t, generated)

	expectedResult := fmt.Sprintf(`unit tests:
  [%s] plus
    @file.10x:8: 3 is not equal to 2

Ran a total of 1 tests
  * 0 succeeded
  * 1 failed
`, codegen_golang.Red("FAILURE"))
	assert.Equal(t, expectedResult, result)
}
