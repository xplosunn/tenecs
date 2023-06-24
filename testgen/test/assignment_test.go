package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestLiteralAssignment(t *testing.T) {
	programString := `package pkg

helloWorld := (): String => {
  result := "hello world!"
  result
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
	formatted := formatter.DisplayModule(*generated)
	assert.Equal(t, expectedOutput, formatted)
}

func TestLiteralRefAssignment(t *testing.T) {
	programString := `package pkg

helloWorld := (): String => {
  result := "hello world!"
  output := result
  output
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
	formatted := formatter.DisplayModule(*generated)
	assert.Equal(t, expectedOutput, formatted)
}

func TestArgAssignment(t *testing.T) {
	programString := `package pkg

strId := (s: String): String => {
  result := s
  result
}
`
	targetFunctionName := "strId"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("foo", testCaseFoo)
  }

  testCaseFoo := (assert: Assert): Void => {
    result := strId("foo")
    expected := "foo"
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

func TestArgRefAssignment(t *testing.T) {
	programString := `package pkg

strId := (s: String): String => {
  result := s
  output := result
  output
}
`
	targetFunctionName := "strId"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("foo", testCaseFoo)
  }

  testCaseFoo := (assert: Assert): Void => {
    result := strId("foo")
    expected := "foo"
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

func TestAssignmentIf(t *testing.T) {
	programString := `package pkg

logPrefix := (isError: Boolean): String => {
  result := if isError {
    "[error]"
  } else {
    "[info]"
  }
  result
}
`
	targetFunctionName := "logPrefix"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("[error]", testCaseError)
    registry.test("[info]", testCaseInfo)
  }

  testCaseError := (assert: Assert): Void => {
    result := logPrefix(true)
    expected := "[error]"
    assert.equal<String>(result, expected)
  }

  testCaseInfo := (assert: Assert): Void => {
    result := logPrefix(false)
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