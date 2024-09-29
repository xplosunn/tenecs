package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
	"testing"
)

func TestFunctionIf(t *testing.T) {
	programString := `
package pkg

filter := (filterFn: (String) -> Boolean, str: String): String => {
  if filterFn(str) {
    str
  } else {
    ""
  }
}
`
	targetFunctionName := "filter"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("foo", testCaseFoo)
  registry.test("", testCase)
})

testCaseFoo := (testkit: UnitTestKit): Void => {
  result := filter(
    (arg0) => {
      true
    },
    "foo"
  )

  expected := "foo"
  testkit.assert.equal<String>(result, expected)
}

testCase := (testkit: UnitTestKit): Void => {
  result := filter(
    (arg0) => {
      false
    },
    "bar"
  )

  expected := ""
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

func TestFunctionWithStdLibInvocation(t *testing.T) {
	programString := `
package pkg

import tenecs.string.join

joinWrapper := (a: String, b: String): String => {
  join(a, b)
}
`
	targetFunctionName := "joinWrapper"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("foobar", testCaseFoobar)
})

testCaseFoobar := (testkit: UnitTestKit): Void => {
  result := joinWrapper("foo", "bar")

  expected := "foobar"
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

func TestFunctionWithList(t *testing.T) {
	programString := `
package pkg

myFunc := (): List<String> => {
  arr := [String]()
  arr
}
`
	targetFunctionName := "myFunc"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("[]", testCase)
})

testCase := (testkit: UnitTestKit): Void => {
  result := myFunc()

  expected := [String]()
  testkit.assert.equal<List<String>>(result, expected)
}
`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, typeErr := typer.TypecheckSingleFile(*parsed)
	if typeErr != nil {
		t.Fatal(type_error.Render(programString, typeErr.(*type_error.TypecheckError)))
	}
	generated, err := testgen.GenerateCached(t, *parsed, *typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := ""
	for _, declaration := range generated {
		formatted += "\n" + formatter.DisplayDeclaration(declaration) + "\n"
	}
	assert.Equal(t, expectedOutput, formatted)
}
