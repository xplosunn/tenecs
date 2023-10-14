package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestSimpleIf(t *testing.T) {
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
	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.Generate(*typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := formatter.DisplayImplementation(*generated)
	assert.Equal(t, expectedOutput, formatted)
}

func TestSequentialIf(t *testing.T) {
	programString := `package pkg

logPrefix := (a: Boolean, isError: Boolean): String => {
  unusedVar := if a {
    "[e]"
  } else {
    "[i]"
  }
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
    registry.test("[error]", testCaseError)
    registry.test("[info]", testCaseInfo)
    registry.test("[error] again", testCaseErroragain)
    registry.test("[info] again", testCaseInfoagain)
  }

  testCaseError := (assert: Assert): Void => {
    result := logPrefix(true, true)

    expected := "[error]"
    assert.equal<String>(result, expected)
  }

  testCaseInfo := (assert: Assert): Void => {
    result := logPrefix(true, false)

    expected := "[info]"
    assert.equal<String>(result, expected)
  }

  testCaseErroragain := (assert: Assert): Void => {
    result := logPrefix(false, true)

    expected := "[error]"
    assert.equal<String>(result, expected)
  }

  testCaseInfoagain := (assert: Assert): Void => {
    result := logPrefix(false, false)

    expected := "[info]"
    assert.equal<String>(result, expected)
  }
}`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.Generate(*typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := formatter.DisplayImplementation(*generated)
	assert.Equal(t, expectedOutput, formatted)
}

func TestThenIf(t *testing.T) {
	programString := `package pkg

logPrefix := (isError: Boolean, isItReally: Boolean): String => {
  if isError {
    if isItReally {
      "[error]"
    } else {
      "[warn]"
    }
  } else {
    "[info]"
  }
}
`
	targetFunctionName := "logPrefix"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("[error]", testCaseError)
    registry.test("[warn]", testCaseWarn)
    registry.test("[info]", testCaseInfo)
  }

  testCaseError := (assert: Assert): Void => {
    result := logPrefix(true, true)

    expected := "[error]"
    assert.equal<String>(result, expected)
  }

  testCaseWarn := (assert: Assert): Void => {
    result := logPrefix(true, false)

    expected := "[warn]"
    assert.equal<String>(result, expected)
  }

  testCaseInfo := (assert: Assert): Void => {
    result := logPrefix(false, true)

    expected := "[info]"
    assert.equal<String>(result, expected)
  }
}`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.Generate(*typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := formatter.DisplayImplementation(*generated)
	assert.Equal(t, expectedOutput, formatted)
}
