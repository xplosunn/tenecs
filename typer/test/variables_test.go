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

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime: Runtime) => {
		output := "Hello world!"
		runtime.console.log(output)
	}
}
`)
	expectedProgram := ast.Program{
		Declarations: []*ast.Declaration{
			{
				VariableType: &types.Void{},
				Name:         "app",
				Expression: &ast.Function{
					VariableType: &types.Function{
						Arguments:  []types.FunctionArgument{},
						ReturnType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
					},
					Block: []ast.Expression{
						ast.Module{
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
										ReturnType: &types.Void{},
									},
									Block: []ast.Expression{
										ast.Declaration{
											VariableType: &types.Void{},
											Name:         "output",
											Expression: ast.Literal{
												VariableType: &types.BasicType{
													Type: "String",
												},
												Literal: parser.LiteralString{
													Value: "\"Hello world!\"",
												},
											},
										},
										ast.WithAccessAndMaybeInvocation{
											VariableType: &types.Void{},
											Over: ast.WithAccessAndMaybeInvocation{
												VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Console"),
												Over: ast.ReferenceAndMaybeInvocation{
													VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
													Name:         "runtime",
												},
												Access: "console",
											},
											Access: "log",
											ArgumentsList: &ast.ArgumentsList{
												Generics: []types.StructFieldVariableType{},
												Arguments: []ast.Expression{
													ast.ReferenceAndMaybeInvocation{
														VariableType: &types.BasicType{
															Type: "String",
														},
														Name: "output",
													},
												},
											},
										},
									},
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
	assert.Equal(t, expectedProgram, program)
}
