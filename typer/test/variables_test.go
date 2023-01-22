package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestMainProgramWithVariable(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
		output := "Hello world!"
		runtime.console.log(output)
	}
}
`)
	expectedProgram := ast.Program{
		Modules: []*ast.Module{
			{
				Name: "app",
				Implements: types.Interface{
					Package: "tenecs.os",
					Name:    "Main",
				},
				Variables: map[string]ast.Expression{
					"main": &ast.Function{
						VariableType: types.Function{
							Arguments: []types.FunctionArgument{
								{
									Name: "runtime",
									VariableType: types.Interface{
										Package: "tenecs.os",
										Name:    "Runtime",
									},
								},
							},
							ReturnType: types.Void{},
						},
						Block: []ast.Expression{
							ast.Declaration{
								VariableType: types.Void{},
								Name:         "output",
								Expression: ast.Literal{
									VariableType: types.BasicType{
										Type: "String",
									},
									Literal: parser.LiteralString{
										Value: "\"Hello world!\"",
									},
								},
							},
							ast.WithAccessAndMaybeInvocation{
								VariableType: types.Void{},
								Over: ast.ReferenceAndMaybeInvocation{
									VariableType: types.Interface{
										Package: "tenecs.os",
										Name:    "Runtime",
									},
									Name: "runtime",
								},
								AccessChain: []ast.AccessAndMaybeInvocation{
									{
										VariableType: types.Interface{
											Package: "tenecs.os",
											Name:    "Console",
										},
										Access: "console",
									},
									{
										VariableType: types.Void{},
										Access:       "log",
										ArgumentsList: &ast.ArgumentsList{
											Arguments: []ast.Expression{
												ast.ReferenceAndMaybeInvocation{
													VariableType: types.BasicType{
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
	}
	assert.Equal(t, expectedProgram, program)
}
