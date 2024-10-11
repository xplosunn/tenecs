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
	programString := `
package pkg

helloWorld := (): String => {
  "hello world!"
}
`
	targetFunctionName := "helloWorld"

	expectedOutput := `
unitTests := UnitTestSuite((registry: UnitTestRegistry): Void => {
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

func TestOneLinerBoolean(t *testing.T) {
	programString := `
package pkg

itIsTrue := (): Boolean => {
  true
}
`
	targetFunctionName := "itIsTrue"

	expectedOutput := `
unitTests := UnitTestSuite((registry: UnitTestRegistry): Void => {
  registry.test("true", testCaseTrue)
})

testCaseTrue := (testkit: UnitTestKit): Void => {
  result := itIsTrue()

  expected := true
  testkit.assert.equal<Boolean>(result, expected)
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
