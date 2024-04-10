package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestMainDirectProgramEmpty(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.os.Main

app := implement Main {
  main := (runtime) => {
		
	}
}
`)
	expectedProgram := ast.Program{
		Package: "main",
		Declarations: []*ast.Declaration{
			{
				Name: "app",
				Expression: ast.Implementation{
					Implements: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
					Variables: map[string]ast.Expression{
						"main": &ast.Function{
							VariableType: &types.Function{
								Arguments: []types.FunctionArgument{
									{
										Name:         "runtime",
										VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
									},
								},
								ReturnType: types.Void(),
							},
							Block: []ast.Expression{},
						},
					},
				},
			},
		},
		StructFunctions:        map[string]*types.Function{},
		NativeFunctions:        map[string]*types.Function{},
		NativeFunctionPackages: map[string]string{},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}

func TestMainProgramEmpty(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.os.Main

app := (): Main => implement Main {
  main := (runtime) => {
		
	}
}
`)
	expectedProgram := ast.Program{
		Package: "main",
		Declarations: []*ast.Declaration{
			{
				Name: "app",
				Expression: &ast.Function{
					VariableType: &types.Function{
						Arguments:  []types.FunctionArgument{},
						ReturnType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
					},
					Block: []ast.Expression{
						ast.Implementation{
							Implements: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
							Variables: map[string]ast.Expression{
								"main": &ast.Function{
									VariableType: &types.Function{
										Arguments: []types.FunctionArgument{
											{
												Name:         "runtime",
												VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
											},
										},
										ReturnType: types.Void(),
									},
									Block: []ast.Expression{},
								},
							},
						},
					},
				},
			},
		},
		StructFunctions:        map[string]*types.Function{},
		NativeFunctions:        map[string]*types.Function{},
		NativeFunctionPackages: map[string]string{},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}

func TestMainProgramReturningStringInBody(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

app := (): Main => implement Main {
  main := (runtime) => {
		"can't return string'"
	}
}
`, "expected type Void but found String")
}

func TestMainProgramMultipleMains(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

app := (): Main => implement Main {
  main := (runtime) => {}
  main := (runtime) => {}
}
`, "duplicate variable 'main'")
}
