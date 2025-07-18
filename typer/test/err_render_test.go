package parser_typer_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
	"testing"
)

func TestSuiteRenderErrorUppercasePackage(t *testing.T) {
	type Case struct {
		program  string
		expected string
	}

	cases := []Case{}

	cases = append(cases, Case{
		program: `package MyPackage`,
		expected: `Error in file file.10x

| 1 | package MyPackage
              ^ package name should start with a lowercase letter
`,
	})

	cases = append(cases, Case{
		program: `package MyPackage

`,
		expected: `Error in file file.10x

| 1 | package MyPackage
              ^ package name should start with a lowercase letter
| 2 | `,
	})

	cases = append(cases, Case{
		program: `
package MyPackage


`,
		expected: `Error in file file.10x
| 1 | 
| 2 | package MyPackage
              ^ package name should start with a lowercase letter
| 3 | 
| 4 | `,
	})

	cases = append(cases, Case{
		program: `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime, anotherRuntime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)
`,
		expected: `Error in file file.10x
| 5  | import tenecs.go.Main
| 6  | 
| 7  | app := Main(
| 8  |   main = (runtime: Runtime, anotherRuntime: Runtime) => {
                ^ expected 1 params but got 2
| 9  |     runtime.console.log("Hello world!")
| 10 |   }`,
	})

	cases = append(cases, Case{
		program: `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    applyToString := (f: (String) ~> Void, strF: () ~> String): Void => {
      f(strF())
    }
    output := (): String => {
      "Hello world!"
    }
    applyToString(runtime.console.log, () => {false})
  }
)
`,
		expected: `Error in file file.10x
| 12 |     output := (): String => {
| 13 |       "Hello world!"
| 14 |     }
| 15 |     applyToString(runtime.console.log, () => {false})
                                                     ^ expected type String but found Boolean
| 16 |   }
| 17 | )`,
	})

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			res, err := parser.ParseString(testCase.program)
			assert.NoError(t, err)

			desugared, err := desugar.Desugar(*res)
			assert.NoError(t, err)

			_, err = typer.TypecheckSingleFile(desugared)
			assert.Error(t, err, "Didn't get an typererror")

			typecheckErr, ok := err.(*type_error.TypecheckError)
			assert.True(t, ok)

			rendered, err := type_error.Render(testCase.program, typecheckErr)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expected, rendered)
		})
	}
}
