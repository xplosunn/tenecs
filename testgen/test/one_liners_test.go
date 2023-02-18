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

func TestOneLinerBoolean(t *testing.T) {
	programString := `package pkg

itIsTrue := (): Boolean => {
  true
}
`
	targetFunctionName := "itIsTrue"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("true", testCasetrue)
  }

  testCasetrue := (assert: Assert): Void => {
    result := module.itIsTrue()
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
	formatted := formatter.DisplayModule(*generated)
	assert.Equal(t, expectedOutput, formatted)
}

func TestOneLinerIf(t *testing.T) {
	programString := `package pkg

logPrefix := (isError: Boolean): String => {
  if isError {
    "[error]"
  } else {
    "[info]"
  }
}
`
	targetFunctionName := "logPrefix"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("[error]", testCaseerror)
    registry.test("[info]", testCaseinfo)
  }

  testCaseerror := (assert: Assert): Void => {
    result := module.logPrefix(true)
    expected := "[error]"
    assert.equal<String>(result, expected)
  }

  testCaseinfo := (assert: Assert): Void => {
    result := module.logPrefix(false)
    expected := "[info]"
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
