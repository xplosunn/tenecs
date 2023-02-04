package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestOneLinerString(t *testing.T) {
	programString := `package pkg

helloWorld := (): String => {
  "hello world!"
}
`
	targetFunctionName := "helloWorld"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("hello world!", testCasehelloworld)
  }

  testCasehelloworld := (assert: Assert): Void => {
    result := module.helloWorld()
    expected := "hello world!"
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
