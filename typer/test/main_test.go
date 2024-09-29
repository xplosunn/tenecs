package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
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
		VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
		Over: ast.Reference{
			VariableType: standard_library.StdLibGetFunctionOrPanic(t, "tenecs.os.Main"),
			PackageName:  ptr("tenecs.os"),
			Name:         "Main",
		},
		Generics: []types.VariableType{},
		Arguments: []ast.Expression{
			&ast.Function{
				VariableType: &types.Function{
					Arguments: []types.FunctionArgument{
						{
							Name:         "runtime",
							VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
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
								Package: "tenecs.os",
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
			Package: "tenecs.os",
			Name:    "Main",
		},
	}
}

func runtimeNativeFunction() *types.Function {
	return &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name: "console",
				VariableType: &types.KnownType{
					Package: "tenecs.os",
					Name:    "Console",
				},
			},
			{
				Name: "http",
				VariableType: &types.KnownType{
					Package:  "tenecs.http",
					Name:     "RuntimeServer",
					Generics: []types.VariableType{},
				},
			},
			{
				Name: "ref",
				VariableType: &types.KnownType{
					Package:  "tenecs.ref",
					Name:     "RefCreator",
					Generics: []types.VariableType{},
				},
			},
		},
		ReturnType: &types.KnownType{
			Package: "tenecs.os",
			Name:    "Runtime",
		},
	}
}

func TestMainDirectProgramEmpty(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.os.Main

app := Main(
  main = (runtime) => {

  }
)
`)
	expectedProgram := ast.Program{
		Package: "main",
		Declarations: []*ast.Declaration{
			{
				Name:       "app",
				Expression: mainWithBlock(t, []ast.Expression{}),
			},
		},
		StructFunctions: map[string]*types.Function{},
		NativeFunctions: map[string]*types.Function{
			"Main": mainNativeFunction(),
		},
		NativeFunctionPackages: map[string]string{
			"Main": "tenecs_os",
		},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}

func TestMainProgramEmpty(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.os.Main

app := Main(
  main = (runtime) => {

  }
)
`)
	expectedProgram := ast.Program{
		Package: "main",
		Declarations: []*ast.Declaration{
			{
				Name:       "app",
				Expression: mainWithBlock(t, []ast.Expression{}),
			},
		},
		StructFunctions: map[string]*types.Function{},
		NativeFunctions: map[string]*types.Function{
			"Main": mainNativeFunction(),
		},
		NativeFunctionPackages: map[string]string{
			"Main": "tenecs_os",
		},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}

func TestMainProgramReturningStringInBody(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

app := Main(
  main = (runtime) => {
    "can't return string'"
  }
)
`, "expected type Void but found String")
}
