package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestExpectedGenericFunctionInvoked4(t *testing.T) {
	program := validProgram(t, testcode.GenericFunctionInvoked4)
	expectedProgram := ast.Program{
		Declarations: []*ast.Declaration{
			{
				Name: "app",
				Expression: &ast.Function{
					VariableType: &types.Function{
						Arguments:  []types.FunctionArgument{},
						ReturnType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
					},
					Block: []ast.Expression{
						ast.Module{
							Implements: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
							Variables: map[string]ast.Expression{
								"identity": &ast.Function{
									VariableType: &types.Function{
										Generics: []string{
											"T",
										},
										Arguments: []types.FunctionArgument{
											{
												Name: "arg",
												VariableType: &types.TypeArgument{
													Name: "T",
												},
											},
										},
										ReturnType: &types.TypeArgument{
											Name: "T",
										},
									},
									Block: []ast.Expression{
										ast.Declaration{
											Name: "result",
											Expression: ast.Reference{
												VariableType: &types.TypeArgument{
													Name: "T",
												},
												Name: "arg",
											},
										},
										ast.Reference{
											VariableType: &types.TypeArgument{
												Name: "T",
											},
											Name: "result",
										},
									},
								},
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
											Name: "output",
											Expression: ast.Literal{
												VariableType: &types.BasicType{
													Type: "String",
												},
												Literal: parser.LiteralString{
													Value: "\"Hello world!\"",
												},
											},
										},
										ast.Declaration{
											Name: "hw",
											Expression: ast.Invocation{
												VariableType: &types.BasicType{
													Type: "String",
												},
												Over: ast.Reference{
													VariableType: &types.Function{
														Arguments: []types.FunctionArgument{
															{
																Name: "arg",
																VariableType: &types.BasicType{
																	Type: "String",
																},
															},
														},
														ReturnType: &types.BasicType{
															Type: "String",
														},
													},
													Name: "identity",
												},
												Generics: []types.StructFieldVariableType{
													&types.BasicType{
														Type: "String",
													},
												},
												Arguments: []ast.Expression{
													ast.Reference{
														VariableType: &types.BasicType{
															Type: "String",
														},
														Name: "output",
													},
												},
											},
										},
										ast.Invocation{
											VariableType: &types.Void{},
											Over: ast.Access{
												VariableType: &types.Function{
													Arguments: []types.FunctionArgument{
														{
															Name: "message",
															VariableType: &types.BasicType{
																Type: "String",
															},
														},
													},
													ReturnType: &types.Void{},
												},
												Over: ast.Access{
													VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Console"),
													Over: ast.Reference{
														VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
														Name:         "runtime",
													},
													Access: "console",
												},
												Access: "log",
											},
											Generics: []types.StructFieldVariableType{},
											Arguments: []ast.Expression{
												ast.Reference{
													VariableType: &types.BasicType{
														Type: "String",
													},
													Name: "hw",
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

func TestExpectedGenericFunctionDoubleInvoked(t *testing.T) {
	program := validProgram(t, testcode.GenericFunctionDoubleInvoked)
	expectedProgram := ast.Program{
		Declarations: []*ast.Declaration{
			{
				Name: "app",
				Expression: &ast.Function{
					VariableType: &types.Function{
						Arguments:  []types.FunctionArgument{},
						ReturnType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
					},
					Block: []ast.Expression{
						ast.Module{
							Implements: standard_library.StdLibGetOrPanic(t, "tenecs.os.Main"),
							Variables: map[string]ast.Expression{
								"identity": &ast.Function{
									VariableType: &types.Function{
										Generics: []string{
											"T",
										},
										Arguments: []types.FunctionArgument{
											{
												Name: "arg",
												VariableType: &types.TypeArgument{
													Name: "T",
												},
											},
										},
										ReturnType: &types.TypeArgument{
											Name: "T",
										},
									},
									Block: []ast.Expression{
										ast.Declaration{
											Name: "output",
											Expression: ast.Invocation{
												VariableType: &types.TypeArgument{
													Name: "T",
												},
												Over: ast.Reference{
													VariableType: &types.Function{
														Arguments: []types.FunctionArgument{
															{
																Name: "arg",
																VariableType: &types.TypeArgument{
																	Name: "T",
																},
															},
														},
														ReturnType: &types.TypeArgument{
															Name: "T",
														},
													},
													Name: "identityFn",
												},
												Generics: []types.StructFieldVariableType{
													&types.TypeArgument{
														Name: "T",
													},
												},
												Arguments: []ast.Expression{
													ast.Reference{
														VariableType: &types.TypeArgument{
															Name: "T",
														},
														Name: "arg",
													},
												},
											},
										},
										ast.Reference{
											VariableType: &types.TypeArgument{
												Name: "T",
											},
											Name: "output",
										},
									},
								},
								"identityFn": &ast.Function{
									VariableType: &types.Function{
										Generics: []string{
											"A",
										},
										Arguments: []types.FunctionArgument{
											{
												Name: "arg",
												VariableType: &types.TypeArgument{
													Name: "A",
												},
											},
										},
										ReturnType: &types.TypeArgument{
											Name: "A",
										},
									},
									Block: []ast.Expression{
										ast.Declaration{
											Name: "result",
											Expression: ast.Reference{
												VariableType: &types.TypeArgument{
													Name: "A",
												},
												Name: "arg",
											},
										},
										ast.Reference{
											VariableType: &types.TypeArgument{
												Name: "A",
											},
											Name: "result",
										},
									},
								},
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
										ast.Invocation{
											VariableType: &types.Void{},
											Over: ast.Access{
												VariableType: &types.Function{
													Arguments: []types.FunctionArgument{
														{
															Name: "message",
															VariableType: &types.BasicType{
																Type: "String",
															},
														},
													},
													ReturnType: &types.Void{},
												},
												Over: ast.Access{
													VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Console"),
													Over: ast.Reference{
														VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
														Name:         "runtime",
													},
													Access: "console",
												},
												Access: "log",
											},
											Generics: []types.StructFieldVariableType{},
											Arguments: []ast.Expression{
												ast.Invocation{
													VariableType: &types.BasicType{
														Type: "String",
													},
													Over: ast.Reference{
														VariableType: &types.Function{
															Arguments: []types.FunctionArgument{
																{
																	Name: "arg",
																	VariableType: &types.BasicType{
																		Type: "String",
																	},
																},
															},
															ReturnType: &types.BasicType{
																Type: "String",
															},
														},
														Name: "identity",
													},
													Generics: []types.StructFieldVariableType{
														&types.BasicType{
															Type: "String",
														},
													},
													Arguments: []ast.Expression{
														ast.Literal{
															VariableType: &types.BasicType{
																Type: "String",
															},
															Literal: parser.LiteralString{
																Value: "\"ciao\"",
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
			},
		},
		StructFunctions:        map[string]*types.Function{},
		NativeFunctions:        map[string]*types.Function{},
		NativeFunctionPackages: map[string]string{},
	}
	assert.Equal(t, expectedProgram, program)
}
