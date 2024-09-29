package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
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
	targetFunctionName := "newPost"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("{title:Breaking news!}", testCaseTitlebreakingnews)
})

testCaseTitlebreakingnews := (testkit: UnitTestKit): Void => {
  result := newPost()

  expected := Post("Breaking news!")
  testkit.assert.equal<Post>(result, expected)
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

func TestStructAccess(t *testing.T) {
	programString := `
package pkg

struct Post(title: String)

postTitle := (post: Post): String => {
  post.title
}
`
	targetFunctionName := "postTitle"

	expectedOutput := `
unitTests := UnitTests((registry: UnitTestRegistry): Void => {
  registry.test("foo", testCaseFoo)
})

testCaseFoo := (testkit: UnitTestKit): Void => {
  result := postTitle(Post("foo"))

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
