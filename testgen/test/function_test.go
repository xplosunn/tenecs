package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestFunctionIf(t *testing.T) {
	programString := `package pkg

filter := (filterFn: (String) -> Boolean, str: String): String => {
  if filterFn(str) {
    str
  } else {
    ""
  }
}
`
	targetFunctionName := "filter"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("foo", testCaseFoo)
    registry.test("", testCase)
  }

  testCaseFoo := (assert: Assert): Void => {
    result := module.filter((arg0) => {
      true
    }, "foo")
    expected := "foo"
    assert.equal<String>(result, expected)
  }

  testCase := (assert: Assert): Void => {
    result := module.filter((arg0) => {
      false
    }, "foo")
    expected := ""
    assert.equal<String>(result, expected)
  }
}`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.Generate(*typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := formatter.DisplayModule(*generated)
	assert.Equal(t, expectedOutput, formatted)
}
