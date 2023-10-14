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
	programString := `package pkg

struct Post(title: String)

newPost := (): Post => {
  Post("Breaking news!")
}
`
	targetFunctionName := "newPost"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("Breaking news!", testCaseBreakingnews)
  }

  testCaseBreakingnews := (assert: Assert): Void => {
    result := newPost()

    expected := Post("Breaking news!")
    assert.equal<Post>(result, expected)
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

func TestStructAccess(t *testing.T) {
	programString := `package pkg

struct Post(title: String)

postTitle := (post: Post): String => {
  post.title
}
`
	targetFunctionName := "postTitle"

	expectedOutput := `implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("foo", testCaseFoo)
  }

  testCaseFoo := (assert: Assert): Void => {
    result := postTitle(Post("foo"))

    expected := "foo"
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
