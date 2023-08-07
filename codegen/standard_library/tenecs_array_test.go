package standard_library_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestFilter(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.array.filter
import tenecs.compare.eq

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("filter", (testkit: UnitTestKit): Void => {
      testkit.assert.equal<Array<String>>([String](), filter<String>([String]("a", "b", "c"), (elem) => false))
      testkit.assert.equal<Array<String>>([String]("a", "b", "c"), filter<String>([String]("a", "b", "c"), (elem) => true))
      testkit.assert.equal<Array<String>>([String]("b"), filter<String>([String]("a", "b", "c"), (a) => eq(a, "b")))
    })
  }
}`
	expectedRunResult := fmt.Sprintf(`myTests:
  [%s] filter

Ran a total of 1 tests
  * 1 succeeded
  * 0 failed
`, codegen.Green("OK"))

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(true, typed)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}
func TestMap(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.array.map
import tenecs.string.join

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("map", (testkit: UnitTestKit): Void => {
      addBang := (s: String): String => { join(s, "!") }
      testkit.assert.equal<Array<String>>([String](), map<String, String>([String](), addBang))
      testkit.assert.equal<Array<String>>([String]("hi!"), map<String, String>([String]("hi"), addBang))
      testkit.assert.equal<Array<String>>([String]("!", "a!", "!", "b!"), map<String, String>([String]("", "a", "", "b"), addBang))
    })
  }
}`
	expectedRunResult := fmt.Sprintf(`myTests:
  [%s] map

Ran a total of 1 tests
  * 1 succeeded
  * 0 failed
`, codegen.Green("OK"))

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
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.array.repeat

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("repeat", (testkit: UnitTestKit): Void => {
      testkit.assert.equal<Array<String>>([String](), repeat<String>("", 0))
      testkit.assert.equal<Array<String>>([String](""), repeat<String>("", 1))
      testkit.assert.equal<Array<String>>([String]("", ""), repeat<String>("", 2))
      testkit.assert.equal<Array<String>>([String]("a"), repeat<String>("a", 1))
      testkit.assert.equal<Array<String>>([String]("a", "a"), repeat<String>("a", 2))
    })
  }
}`
	expectedRunResult := fmt.Sprintf(`myTests:
  [%s] repeat

Ran a total of 1 tests
  * 1 succeeded
  * 0 failed
`, codegen.Green("OK"))

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
