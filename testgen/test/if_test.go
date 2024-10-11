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
	programString := `
package pkg

logPrefix := (isError: Boolean): String => {
  if isError {
    "[error]"
  } else {
    "[info]"
  }
}
`
	targetFunctionName := "logPrefix"

	expectedOutput := `
unitTests := UnitTestSuite((registry: UnitTestRegistry): Void => {
  registry.test("[error]", testCaseError)
  registry.test("[info]", testCaseInfo)
})

testCaseError := (testkit: UnitTestKit): Void => {
  result := logPrefix(true)

  expected := "[error]"
  testkit.assert.equal<String>(result, expected)
}

testCaseInfo := (testkit: UnitTestKit): Void => {
  result := logPrefix(false)

  expected := "[info]"
  testkit.assert.equal<String>(result, expected)
}
`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.GenerateCached(t, *parsed, *typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := ""
	for _, declaration := range generated {
		formatted += "\n" + formatter.DisplayDeclaration(declaration) + "\n"
	}
	assert.Equal(t, expectedOutput, formatted)
}

func TestSequentialIf(t *testing.T) {
	programString := `
package pkg

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

	expectedOutput := `
unitTests := UnitTestSuite((registry: UnitTestRegistry): Void => {
  registry.test("[error]", testCaseError)
  registry.test("[info]", testCaseInfo)
  registry.test("[error] again", testCaseErroragain)
  registry.test("[info] again", testCaseInfoagain)
})

testCaseError := (testkit: UnitTestKit): Void => {
  result := logPrefix(true, true)

  expected := "[error]"
  testkit.assert.equal<String>(result, expected)
}

testCaseInfo := (testkit: UnitTestKit): Void => {
  result := logPrefix(true, false)

  expected := "[info]"
  testkit.assert.equal<String>(result, expected)
}

testCaseErroragain := (testkit: UnitTestKit): Void => {
  result := logPrefix(false, true)

  expected := "[error]"
  testkit.assert.equal<String>(result, expected)
}

testCaseInfoagain := (testkit: UnitTestKit): Void => {
  result := logPrefix(false, false)

  expected := "[info]"
  testkit.assert.equal<String>(result, expected)
}
`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.GenerateCached(t, *parsed, *typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := ""
	for _, declaration := range generated {
		formatted += "\n" + formatter.DisplayDeclaration(declaration) + "\n"
	}
	assert.Equal(t, expectedOutput, formatted)
}

func TestThenIf(t *testing.T) {
	programString := `
package pkg

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

	expectedOutput := `
unitTests := UnitTestSuite((registry: UnitTestRegistry): Void => {
  registry.test("[error]", testCaseError)
  registry.test("[warn]", testCaseWarn)
  registry.test("[info]", testCaseInfo)
})

testCaseError := (testkit: UnitTestKit): Void => {
  result := logPrefix(true, true)

  expected := "[error]"
  testkit.assert.equal<String>(result, expected)
}

testCaseWarn := (testkit: UnitTestKit): Void => {
  result := logPrefix(true, false)

  expected := "[warn]"
  testkit.assert.equal<String>(result, expected)
}

testCaseInfo := (testkit: UnitTestKit): Void => {
  result := logPrefix(false, true)

  expected := "[info]"
  testkit.assert.equal<String>(result, expected)
}
`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.GenerateCached(t, *parsed, *typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := ""
	for _, declaration := range generated {
		formatted += "\n" + formatter.DisplayDeclaration(declaration) + "\n"
	}
	assert.Equal(t, expectedOutput, formatted)
}
