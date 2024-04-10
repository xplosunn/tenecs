package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"strconv"
	"testing"
)

func TestMultipleFilesInterfaceWithSeparateImplementationVariableString(t *testing.T) {
	validProgramFromSinglePackage(t, []string{
		`
package main

interface A {
  a: () -> String
}
`, `
package main

app := (): A => implement A {
  a := () => ""
}
`,
	})
}
func TestMultipleFilesInterfaceReturningAnotherInterfaceInVariable(t *testing.T) {
	f1 := `
package main

interface Goods {
  name: () -> String
}
`
	f2 := `
package main

interface Factory {
  produce: () -> Goods
}
`
	validProgramFromSinglePackage(t, []string{f1, f2})
	validProgramFromSinglePackage(t, []string{f2, f1})
}

func TestMultipleFilesDuplicateStruct(t *testing.T) {
	invalidProgramFromSinglePackage(t, []string{
		`
package main

struct Dup(a: String)
`, `
package main

struct Dup(a: String)
`,
	}, "type already exists Dup")
}

func TestMultipleFilesDuplicateInterface(t *testing.T) {
	invalidProgramFromSinglePackage(t, []string{
		`
package main

interface A {}
`, `
package main

interface A {}
`,
	}, "type already exists A")
}

func validProgramFromSinglePackage(t *testing.T, fileContents []string) ast.Program {
	assert.NotZero(t, fileContents)
	parsedFiles := map[string]parser.FileTopLevel{}
	pkgName := ""
	for i, content := range fileContents {
		res, err := parser.ParseString(content)
		assert.NoError(t, err)
		parsedFiles["f"+strconv.Itoa(i)+".10x"] = *res
		if pkgName == "" {
			for i, name := range res.Package.DotSeparatedNames {
				if i > 0 {
					pkgName += "."
				}
				pkgName += name.String
			}
		}
	}

	p, typeErr := typer.TypecheckPackage(pkgName, parsedFiles)
	if typeErr != nil {
		//TODO re-add:
		//t.Fatal(type_error.Render(program, typeErr.(*type_error.TypecheckError)))
		t.Fatal(typeErr.Error())
	}
	return *p
}

func invalidProgramFromSinglePackage(t *testing.T, fileContents []string, errorMessage string) {
	assert.NotZero(t, fileContents)
	parsedFiles := map[string]parser.FileTopLevel{}
	pkgName := ""
	for i, content := range fileContents {
		res, err := parser.ParseString(content)
		assert.NoError(t, err)
		parsedFiles["f"+strconv.Itoa(i)+".10x"] = *res
		if pkgName == "" {
			for i, name := range res.Package.DotSeparatedNames {
				if i > 0 {
					pkgName += "."
				}
				pkgName += name.String
			}
		}
	}

	_, typeErr := typer.TypecheckPackage(pkgName, parsedFiles)
	assert.Error(t, typeErr, "Didn't get an typererror")
	assert.Equal(t, errorMessage, typeErr.Error())
}
