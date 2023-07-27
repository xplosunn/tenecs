package standard_library_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"os"
	"os/exec"
	"path/filepath"
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

func TestRepeat(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestRegistry
import tenecs.test.Assert
import tenecs.array.repeat

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("repeat", (assert: Assert): Void => {
      assert.equal<Array<String>>([String](), repeat<String>("", 0))
      assert.equal<Array<String>>([String](""), repeat<String>("", 1))
      assert.equal<Array<String>>([String]("", ""), repeat<String>("", 2))
      assert.equal<Array<String>>([String]("a"), repeat<String>("a", 1))
      assert.equal<Array<String>>([String]("a", "a"), repeat<String>("a", 2))
    })
  }
}`
	expectedRunResult := `myTests:
  [OK] repeat
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(true, typed)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func createFileAndRun(t *testing.T, fileContent string) string {
	dir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	filePath := filepath.Join(dir, t.Name()+".go")

	_, err = os.Create(filePath)

	contentBytes := []byte(fileContent)
	err = os.WriteFile(filePath, contentBytes, 0644)
	assert.NoError(t, err)

	cmd := exec.Command("go", "run", filePath)
	cmd.Dir = dir
	outputBytes, err := cmd.Output()
	t.Log(dir)
	assert.NoError(t, err)
	return string(outputBytes)
}
