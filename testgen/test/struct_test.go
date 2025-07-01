package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"testing"
)

func TestStructInstance(t *testing.T) {
	programString := `
package pkg

struct Post(title: String)

newPost := (): Post => {
  Post("Breaking news!")
}
`
	targetFunctionName := ast.Ref{
		Package: "pkg",
		Name:    "newPost",
	}

	expectedOutput := `
_ := UnitTest("{title:Breaking news!}", (testkit: UnitTestKit): Void => {
  result := newPost()

  expected := Post("Breaking news!")
  testkit.assert.equal<Post>(result, expected)
})
`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	desugared := desugar.Desugar(*parsed)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)
	generated, err := testgen.GenerateCached(t, *parsed, *typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := ""
	for _, declaration := range generated {
		formatted += "\n" + formatter.DisplayDeclaration(declaration) + "\n"
	}
	assert.Equal(t, expectedOutput, formatted)
}

func TestStructAccess(t *testing.T) {
	programString := `
package pkg

struct Post(title: String)

postTitle := (post: Post): String => {
  post.title
}
`
	targetFunctionName := ast.Ref{
		Package: "pkg",
		Name:    "postTitle",
	}

	expectedOutput := `
_ := UnitTest("foo", (testkit: UnitTestKit): Void => {
  result := postTitle(Post("foo"))

  expected := "foo"
  testkit.assert.equal<String>(result, expected)
})
`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	desugared := desugar.Desugar(*parsed)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)
	generated, err := testgen.GenerateCached(t, *parsed, *typed, targetFunctionName)
	assert.NoError(t, err)
	formatted := ""
	for _, declaration := range generated {
		formatted += "\n" + formatter.DisplayDeclaration(declaration) + "\n"
	}
	assert.Equal(t, expectedOutput, formatted)
}
