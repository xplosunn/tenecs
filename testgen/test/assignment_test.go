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
	programString := `
package pkg

helloWorld := (): String => {
  result := "hello world!"
  result
}
`
	targetFunctionName := "helloWorld"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("hello world!", testCaseHelloworld)
})

testCaseHelloworld := (testkit: UnitTestKit): Void => {
  result := helloWorld()

  expected := "hello world!"
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

func TestLiteralRefAssignment(t *testing.T) {
	programString := `
package pkg

helloWorld := (): String => {
  result := "hello world!"
  output := result
  output
}
`
	targetFunctionName := "helloWorld"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("hello world!", testCaseHelloworld)
})

testCaseHelloworld := (testkit: UnitTestKit): Void => {
  result := helloWorld()

  expected := "hello world!"
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

func TestArgAssignment(t *testing.T) {
	programString := `
package pkg

strId := (s: String): String => {
  result := s
  result
}
`
	targetFunctionName := "strId"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("foo", testCaseFoo)
})

testCaseFoo := (testkit: UnitTestKit): Void => {
  result := strId("foo")

  expected := "foo"
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

func TestArgRefAssignment(t *testing.T) {
	programString := `
package pkg

strId := (s: String): String => {
  result := s
  output := result
  output
}
`
	targetFunctionName := "strId"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("foo", testCaseFoo)
})

testCaseFoo := (testkit: UnitTestKit): Void => {
  result := strId("foo")

  expected := "foo"
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

func TestAssignmentIf(t *testing.T) {
	programString := `
package pkg

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

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
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
