package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func ptr[T any](t T) *T {
	return &t
}

func mainWithBlock(t *testing.T, block []ast.Expression) ast.Invocation {
	return ast.Invocation{
		VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.go.Main"),
		Over: ast.Reference{
			VariableType: standard_library.StdLibGetFunctionOrPanic(t, "tenecs.go.Main"),
			PackageName:  ptr("tenecs.go"),
			Name:         "Main",
		},
		Generics: []types.VariableType{},
		Arguments: []ast.Expression{
			&ast.Function{
				VariableType: &types.Function{
					Arguments: []types.FunctionArgument{
						{
							Name:         "runtime",
							VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.go.Runtime"),
						},
					},
					ReturnType: types.Void(),
				},
				Block: block,
			},
		},
	}
}

func mainNativeFunction() *types.Function {
	return &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name: "main",
				VariableType: &types.Function{
					Arguments: []types.FunctionArgument{
						{
							Name: "runtime",
							VariableType: &types.KnownType{
								Package: "tenecs.go",
								Name:    "Runtime",
							},
						},
					},
					ReturnType: &types.KnownType{
						Name: "Void",
					},
				},
			},
		},
		ReturnType: &types.KnownType{
			Package: "tenecs.go",
			Name:    "Main",
		},
	}
}

func TestMainDirectProgramEmpty(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.go.Main

app := Main(
  main = (runtime) => {
    null
  }
)
`)
	expectedProgram := ast.Program{
		Declarations: map[ast.Ref]ast.Expression{
			ast.Ref{
				Package: "main",
				Name:    "app",
			}: mainWithBlock(t, []ast.Expression{
				ast.Literal{
					VariableType: &types.KnownType{
						Name: "Void",
					},
					Literal: parser.LiteralNull{
						Value: true,
					},
				},
			}),
		},
		TypeAliases:     map[ast.Ref]ast.TypeAlias{},
		StructFunctions: map[ast.Ref]*types.Function{},
		NativeFunctions: map[ast.Ref]*types.Function{
			ast.Ref{
				Package: "tenecs_go",
				Name:    "Main",
			}: mainNativeFunction(),
		},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, ast.EmptyCodePoints(program))
}

func TestMainProgramReturningStringInBody(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Main

app := Main(
  main = (runtime) => {
    "can't return string'"
  }
)
`, "expected type Void but found String")
}
