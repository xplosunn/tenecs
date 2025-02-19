package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestMainProgramWithVariable(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    output := "Hello world!"
    runtime.console.log(output)
  }
)
`)
	expectedProgram := ast.Program{
		Declarations: map[ast.Ref]ast.Expression{
			ast.Ref{
				Package: "main",
				Name:    "app",
			}: mainWithBlock(t, []ast.Expression{
				ast.Declaration{
					Name: "output",
					Expression: ast.Literal{
						VariableType: types.String(),
						Literal: parser.LiteralString{
							Value: "\"Hello world!\"",
						},
					},
				},
				ast.Invocation{
					VariableType: types.Void(),
					Over: ast.Access{
						VariableType: &types.Function{
							Arguments: []types.FunctionArgument{
								{
									Name:         "message",
									VariableType: types.String(),
								},
							},
							ReturnType: types.Void(),
						},
						Over: ast.Access{
							VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.go.Console"),
							Over: ast.Reference{
								VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.go.Runtime"),
								Name:         "runtime",
							},
							Access: "console",
						},
						Access: "log",
					},
					Generics: []types.VariableType{},
					Arguments: []ast.Expression{
						ast.Reference{
							VariableType: types.String(),
							Name:         "output",
						},
					},
				},
			}),
		},
		StructFunctions: map[ast.Ref]*types.Function{},
		NativeFunctions: map[ast.Ref]*types.Function{
			ast.Ref{
				Package: "tenecs_go",
				Name:    "Main",
			}: mainNativeFunction(),
			ast.Ref{
				Package: "tenecs_go",
				Name:    "Runtime",
			}: runtimeNativeFunction(),
		},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}

func TestInvalidVariableName(t *testing.T) {
	program := `package pk

true := false
`

	invalidProgram(t, program, "Variable can't be named 'true'")
}

func TestInvalidLocalVariableName(t *testing.T) {
	program := `package pk

_ := (): Void => {
  false := true
}
`

	invalidProgram(t, program, "Variable can't be named 'false'")
}
