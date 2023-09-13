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
    registry.test("hello world!", testCaseHelloworld)
  }

  testCaseHelloworld := (assert: Assert): Void => {
    result := helloWorld()

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
	formatted := formatter.DisplayImplementation(*generated)
	assert.Equal(t, expectedOutput, formatted)
}

func TestOneLinerBoolean(t *testing.T) {
	programString := `package pkg

itIsTrue := (): Boolean => {
  true
}
`
	targetFunctionName := "itIsTrue"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("true", testCaseTrue)
  }

  testCaseTrue := (assert: Assert): Void => {
    result := itIsTrue()

    expected := true
    assert.equal<Boolean>(result, expected)
  }
}`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.Generate(*typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := formatter.DisplayImplementation(*generated)
	assert.Equal(t, expectedOutput, formatted)
}
