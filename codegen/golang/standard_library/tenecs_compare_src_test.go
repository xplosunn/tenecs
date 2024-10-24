package standard_library_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	golang2 "github.com/xplosunn/tenecs/codegen/golang"
	"github.com/xplosunn/tenecs/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestEq(t *testing.T) {
	program := `package test

import tenecs.test.UnitTest
import tenecs.test.UnitTestKit
import tenecs.compare.eq

_ := UnitTest("eq", (testkit: UnitTestKit): Void => {
  testkit.assert.equal(true, eq(true, true))
  testkit.assert.equal(false, eq(true, false))
  testkit.assert.equal(true, eq("", ""))
  testkit.assert.equal(false, eq("a", "b"))
})`
	expectedRunResult := fmt.Sprintf(`unit tests:
  [%s] eq

Ran a total of 1 tests
  * 1 succeeded
  * 0 failed
`, golang2.Green("OK"))

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := golang2.GenerateProgramTest(typed)

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}
