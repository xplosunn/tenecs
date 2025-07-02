package codegen_golang_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"testing"
)

var runtime = `func runtime() tenecs_go_Runtime {
	return tenecs_go_Runtime{
		_console: tenecs_go_Console{
			_log: func(Pmessage any) any {
				fmt.Println(Pmessage)
				return nil
			},
		},
		_ref: tenecs_ref_RefCreator{
			_new: func(Pvalue any) any {
				var ref any = Pvalue
				return tenecs_ref_Ref{
					_get: func() any {
						return ref
					},
					_set: func(value any) any {
						ref = value
						return nil
					},
					_modify: func(f any) any {
						ref = f.(func(any) any)(ref)
						return nil
					},
				}

				return nil
			},
		},
	}
}
`

func TestGenerateAndRunTest(t *testing.T) {
	program := `package test

import tenecs.test.UnitTest
import tenecs.test.UnitTestSuite
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry

helloWorld := (): String => {
  "hello world!"
}

_ := UnitTestSuite(
  "My Tests",
  tests = (registry: UnitTestRegistry): Void => {
    registry.test("hello world function", testCaseHelloworld)
  }
)

_ := UnitTest("unitHello", testCaseHelloworld)

testCaseHelloworld := (testkit: UnitTestKit): Void => {
  result := helloWorld()
  expected := "hello world!"
  testkit.assert.equal<String>(result, expected)
}
`
	expectedRunResult := fmt.Sprintf(`unit tests:
  [%s] unitHello
My Tests:
  [%s] hello world function

Ran a total of 2 tests
  * 2 succeeded
  * 0 failed
`, codegen_golang.Green("OK"), codegen_golang.Green("OK"))

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramTest(typed, codegen.FindTests(typed))
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}
func TestGenerateAndRunTestWithManyTests(t *testing.T) {
	program := `package test

import tenecs.test.UnitTest
import tenecs.test.UnitTestSuite
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry

helloWorld := (): String => {
  "hello world!"
}

_ := UnitTest("unitHello1", testCaseHelloworld)
_ := UnitTest("unitHello2", testCaseHelloworld)
_ := UnitTest("unitHello3", testCaseHelloworld)
_ := UnitTest("unitHello4", testCaseHelloworld)
_ := UnitTest("unitHello5", testCaseHelloworld)
_ := UnitTest("unitHello6", testCaseHelloworld)
_ := UnitTest("unitHello7", testCaseHelloworld)
_ := UnitTest("unitHello8", testCaseHelloworld)
_ := UnitTest("unitHello9", testCaseHelloworld)

testCaseHelloworld := (testkit: UnitTestKit): Void => {
  result := helloWorld()
  expected := "hello world!"
  testkit.assert.equal<String>(result, expected)
}
`

	greenOk := codegen_golang.Green("OK")
	expectedRunResult := fmt.Sprintf(`unit tests:
  [%s] unitHello1
  [%s] unitHello2
  [%s] unitHello3
  [%s] unitHello4
  [%s] unitHello5
  [%s] unitHello6
  [%s] unitHello7
  [%s] unitHello8
  [%s] unitHello9

Ran a total of 9 tests
  * 9 succeeded
  * 0 failed
`, greenOk, greenOk, greenOk, greenOk, greenOk, greenOk, greenOk, greenOk, greenOk)

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramTest(typed, codegen.FindTests(typed))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithStandardLibraryFunction(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main
import tenecs.string.join

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log(join("Hello ", "world!"))
  }
)`

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, ast.Ref{
		Package: "main",
		Name:    "app",
	})
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithImportAlias(t *testing.T) {
	program := testcode.ImportAliasMain

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, ast.Ref{
		Package: "main",
		Name:    "app",
	})
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithStruct(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

struct Post(title: String)

app := Main(
  main = (runtime: Runtime) => {
    post := Post("the title")
    runtime.console.log(post.title)
  }
)`

	expectedRunResult := "the title\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, ast.Ref{
		Package: "main",
		Name:    "app",
	})
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMain(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)`

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, ast.Ref{
		Package: "main",
		Name:    "app",
	})
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithRecursion(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main
import tenecs.int.times
import tenecs.int.minus
import tenecs.compare.eq
import tenecs.json.jsonInt

factorial := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    times(i, factorial(minus(i, 1)))
  }
}

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log(jsonInt().toJson(factorial(5)))
  }
)`

	expectedRunResult := "120\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, ast.Ref{
		Package: "main",
		Name:    "app",
	})
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithWhen(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main
import tenecs.json.jsonInt
import tenecs.string.join

struct Post(title: String)

struct BlogPost(title: String)

toString := (input: Int | String | Post | BlogPost): String => {
  when input {
    is i: Int => {
      jsonInt().toJson(i)
    }
    is s: String => {
      s
    }
    is p: Post => {
      join("post:", p.title)
    }
    is b: BlogPost => {
      join("blogpost:", b.title)
    }
  }
}

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log(toString("is it 10?"))
    runtime.console.log(toString(10))
    runtime.console.log(toString(Post("wee")))
    runtime.console.log(toString(BlogPost("wee2")))
  }
)`

	expectedRunResult := `is it 10?
10
post:wee
blogpost:wee2
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, ast.Ref{
		Package: "main",
		Name:    "app",
	})
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateShortCircuitTwice(t *testing.T) {
	program := testcode.ShortCircuitTwice

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared, err := desugar.Desugar(*parsed)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramNonRunnable(typed)
	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))
}
