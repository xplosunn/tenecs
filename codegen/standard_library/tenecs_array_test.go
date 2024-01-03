package standard_library_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
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

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramTest(typed)

	output, err := golang.RunCodeBlockingAndReturningOutputWhenFinished(generated)
	assert.NoError(t, err)
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

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramTest(typed)

	output, err := golang.RunCodeBlockingAndReturningOutputWhenFinished(generated)
	assert.NoError(t, err)
	assert.Equal(t, expectedRunResult, output)
}
func TestFlatMap(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.array.flatMap
import tenecs.string.join

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("flatMap", (testkit: UnitTestKit): Void => {
      addBang := (s: String): Array<String> => { [](s, "!") }
      testkit.assert.equal<Array<String>>([String](), flatMap<String, String>([String](), addBang))
      testkit.assert.equal<Array<String>>([String]("hi", "!"), flatMap<String, String>([String]("hi"), addBang))
      testkit.assert.equal<Array<String>>([String]("a", "!", "b", "!"), flatMap<String, String>([String]("a", "b"), addBang))
    })
  }
}`
	expectedRunResult := fmt.Sprintf(`myTests:
  [%s] flatMap

Ran a total of 1 tests
  * 1 succeeded
  * 0 failed
`, codegen.Green("OK"))

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramTest(typed)

	output, err := golang.RunCodeBlockingAndReturningOutputWhenFinished(generated)
	assert.NoError(t, err)
	assert.Equal(t, expectedRunResult, output)
}

func TestFold(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.array.fold
import tenecs.string.join

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("fold", (testkit: UnitTestKit): Void => {
      testkit.assert.equal<String>("r", fold<Boolean, String>([Boolean](), "r", (acc, elem) => { join(acc, "!") }))
      testkit.assert.equal<String>("_ab", fold<String, String>([]("a", "b"), "_", (acc, elem) => { join(acc, elem) }))
    })
  }
}`
	expectedRunResult := fmt.Sprintf(`myTests:
  [%s] fold

Ran a total of 1 tests
  * 1 succeeded
  * 0 failed
`, codegen.Green("OK"))

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramTest(typed)

	output, err := golang.RunCodeBlockingAndReturningOutputWhenFinished(generated)
	assert.NoError(t, err)
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

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramTest(typed)

	output, err := golang.RunCodeBlockingAndReturningOutputWhenFinished(generated)
	assert.NoError(t, err)
	assert.Equal(t, expectedRunResult, output)
}
