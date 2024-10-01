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
		Package: "main",
		Declarations: []*ast.Declaration{
			{
				Name: "app",
				Expression: mainWithBlock(t, []ast.Expression{
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
		},
		StructFunctions: map[string]*types.Function{},
		NativeFunctions: map[string]*types.Function{
			"Main":    mainNativeFunction(),
			"Runtime": runtimeNativeFunction(),
		},
		NativeFunctionPackages: map[string]string{
			"Main":    "tenecs_go",
			"Runtime": "tenecs_go",
		},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}
