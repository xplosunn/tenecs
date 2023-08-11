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
    result := filter(
      (arg0) => {
        true
      },
      "foo"
    )

    expected := "foo"
    assert.equal<String>(result, expected)
  }

  testCase := (assert: Assert): Void => {
    result := filter(
      (arg0) => {
        false
      },
      "bar"
    )

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

func TestFunctionWithStdLibInvocation(t *testing.T) {
	programString := `package pkg

import tenecs.string.join

joinWrapper := (a: String, b: String): String => {
  join(a, b)
}
`
	targetFunctionName := "joinWrapper"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("foobar", testCaseFoobar)
  }

  testCaseFoobar := (assert: Assert): Void => {
    result := joinWrapper("foo", "bar")

    expected := "foobar"
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

func TestFunctionWithArray(t *testing.T) {
	programString := `package pkg

myFunc := (): Array<String> => {
  arr := [String]()
  arr
}
`
	targetFunctionName := "myFunc"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("[]", testCase)
  }

  testCase := (assert: Assert): Void => {
    result := myFunc()

    expected := [String]()
    assert.equal<Array<String>>(result, expected)
  }
}`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, typeErr := typer.Typecheck(*parsed)
	if typeErr != nil {
		t.Fatal(type_error.Render(programString, typeErr.(*type_error.TypecheckError)))
	}
	generated, err := testgen.Generate(*typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := formatter.DisplayModule(*generated)
	assert.Equal(t, expectedOutput, formatted)
}
