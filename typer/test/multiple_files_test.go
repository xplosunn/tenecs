package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"strconv"
	"testing"
)

func TestMultipleFilesWithSameImport(t *testing.T) {
	validProgramFromFileContents(t, []string{
		`
package main

import tenecs.error.Error

errNotFound := (): Error => {
  Error("not found")
}
`, `
package main

import tenecs.error.Error

errInvalid := (): Error => {
  Error("invalid")
}
`,
	})
}

func TestMultipleFilesStructWithSeparateImplementationVariableString(t *testing.T) {
	validProgramFromFileContents(t, []string{
		`
package main

struct A(
  a: () ~> String
)
`, `
package main

app := (): A => A(
  a = () => ""
)
`,
	})
}

func TestMultipleFilesStructReturningAnotherInterfaceInVariable(t *testing.T) {
	f1 := `
package main

struct Goods(
  name: () ~> String
)
`
	f2 := `
package main

struct Factory(
  produce: () ~> Goods
)
`
	validProgramFromFileContents(t, []string{f1, f2})
	validProgramFromFileContents(t, []string{f2, f1})
}

func TestMultipleFilesDuplicateStruct(t *testing.T) {
	invalidProgramFromFileContents(t, []string{
		`
package main

struct Dup(a: String)
`, `
package main

struct Dup(a: String)
`,
	}, "type already exists: main.Dup")
}

func TestMultiplePackagesStructImport(t *testing.T) {
	f1 := `
package colors

struct Red()
`
	f2 := `
package main

import colors.Red

struct RedWrapper(red: Red)
`
	validProgramFromFileContents(t, []string{f1, f2})
	validProgramFromFileContents(t, []string{f2, f1})
}

func validProgramFromFileContents(t *testing.T, fileContents []string) ast.Program {
	assert.NotZero(t, fileContents)
	parsedFiles := map[string]parser.FileTopLevel{}
	for i, content := range fileContents {
		res, err := parser.ParseString(content)
		assert.NoError(t, err)
		parsedFiles["f"+strconv.Itoa(i)+".10x"] = *res
	}

	p, typeErr := typer.TypecheckPackages(parsedFiles)
	if typeErr != nil {
		//TODO re-add:
		//t.Fatal(type_error.Render(program, typeErr.(*type_error.TypecheckError)))
		t.Fatal(typeErr.Error())
	}
	return *p
}

func invalidProgramFromFileContents(t *testing.T, fileContents []string, errorMessage string) {
	assert.NotZero(t, fileContents)
	parsedFiles := map[string]parser.FileTopLevel{}
	for i, content := range fileContents {
		res, err := parser.ParseString(content)
		assert.NoError(t, err)
		parsedFiles["f"+strconv.Itoa(i)+".10x"] = *res
	}

	_, typeErr := typer.TypecheckPackages(parsedFiles)
	assert.Error(t, typeErr, "Didn't get an typererror")
	assert.Equal(t, errorMessage, typeErr.Error())
}
