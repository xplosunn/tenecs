package standard_library_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestMap(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestRegistry
import tenecs.test.Assert
import tenecs.array.map
import tenecs.string.join

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("map", (assert: Assert): Void => {
      addBang := (s: String): String => { join(s, "!") }
      assert.equal<Array<String>>([String](), map<String, String>([String](), addBang))
      assert.equal<Array<String>>([String]("hi!"), map<String, String>([String]("hi"), addBang))
      assert.equal<Array<String>>([String]("!", "a!", "!", "b!"), map<String, String>([String]("", "a", "", "b"), addBang))
    })
  }
}`
	expectedRunResult := `myTests:
  [OK] map
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(true, typed)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}
