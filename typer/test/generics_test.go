package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestGenericFunctionDeclared(t *testing.T) {
	validProgram(t, testcode.GenericFunctionDeclared)
}

func TestGenericFunctionInvoked(t *testing.T) {
	validProgram(t, testcode.GenericFunctionInvoked1)
	validProgram(t, testcode.GenericFunctionInvoked2)
	validProgram(t, testcode.GenericFunctionInvoked3)
	program := validProgram(t, testcode.GenericFunctionInvoked4)
	expectedProgram := ast.Program{
		Declarations: []*ast.Declaration{
			{
				VariableType: types.BasicType{
					Type: "Void",
				},
				Name: "app",
				Expression: &ast.Function{
					VariableType: types.Function{
						Arguments: []types.FunctionArgument{},
						ReturnType: types.Interface{
							Package: "tenecs.os",
							Name:    "Main",
						},
					},
					Block: []ast.Expression{
						ast.Module{
							Implements: types.Interface{
								Package: "tenecs.os",
								Name:    "Main",
							},
							Variables: map[string]ast.Expression{
								"identity": ast.Function{
									VariableType: types.Function{
										Generics: []string{
											"T",
										},
										Arguments: []types.FunctionArgument{
											{
												Name: "arg",
												VariableType: types.TypeArgument{
													Name: "T",
												},
											},
										},
										ReturnType: types.TypeArgument{
											Name: "T",
										},
									},
									Block: []ast.Expression{
										ast.Declaration{
											VariableType: types.Void{},
											Name:         "result",
											Expression: ast.ReferenceAndMaybeInvocation{
												VariableType: types.TypeArgument{
													Name: "T",
												},
												Name: "arg",
											},
										},
										ast.ReferenceAndMaybeInvocation{
											VariableType: types.TypeArgument{
												Name: "T",
											},
											Name: "result",
										},
									},
								},
								"main": ast.Function{
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
										ast.Declaration{
											VariableType: types.Void{},
											Name:         "hw",
											Expression: ast.ReferenceAndMaybeInvocation{
												VariableType: types.BasicType{
													Type: "String",
												},
												Name: "identity",
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
				},
			},
		},
		StructFunctions: map[string]types.Function{},
	}
	assert.Equal(t, expectedProgram, program)
}

func TestGenericFunctionDoubleInvoked(t *testing.T) {
	program := validProgram(t, testcode.GenericFunctionDoubleInvoked)
	expectedProgram := ast.Program{
		Declarations: []*ast.Declaration{
			{
				VariableType: types.BasicType{
					Type: "Void",
				},
				Name: "app",
				Expression: &ast.Function{
					VariableType: types.Function{
						Arguments: []types.FunctionArgument{},
						ReturnType: types.Interface{
							Package: "tenecs.os",
							Name:    "Main",
						},
					},
					Block: []ast.Expression{
						ast.Module{
							Implements: types.Interface{
								Package: "tenecs.os",
								Name:    "Main",
							},
							Variables: map[string]ast.Expression{
								"identity": ast.Function{
									VariableType: types.Function{
										Generics: []string{
											"T",
										},
										Arguments: []types.FunctionArgument{
											{
												Name: "arg",
												VariableType: types.TypeArgument{
													Name: "T",
												},
											},
										},
										ReturnType: types.TypeArgument{
											Name: "T",
										},
									},
									Block: []ast.Expression{
										ast.Declaration{
											VariableType: types.Void{},
											Name:         "output",
											Expression: ast.ReferenceAndMaybeInvocation{
												VariableType: types.TypeArgument{
													Name: "T",
												},
												Name: "identityFn",
												ArgumentsList: &ast.ArgumentsList{
													Arguments: []ast.Expression{
														ast.ReferenceAndMaybeInvocation{
															VariableType: types.TypeArgument{
																Name: "T",
															},
															Name: "arg",
														},
													},
												},
											},
										},
										ast.ReferenceAndMaybeInvocation{
											VariableType: types.TypeArgument{
												Name: "T",
											},
											Name: "output",
										},
									},
								},
								"identityFn": ast.Function{
									VariableType: types.Function{
										Generics: []string{
											"A",
										},
										Arguments: []types.FunctionArgument{
											{
												Name: "arg",
												VariableType: types.TypeArgument{
													Name: "A",
												},
											},
										},
										ReturnType: types.TypeArgument{
											Name: "A",
										},
									},
									Block: []ast.Expression{
										ast.Declaration{
											VariableType: types.Void{},
											Name:         "result",
											Expression: ast.ReferenceAndMaybeInvocation{
												VariableType: types.TypeArgument{
													Name: "A",
												},
												Name: "arg",
											},
										},
										ast.ReferenceAndMaybeInvocation{
											VariableType: types.TypeArgument{
												Name: "A",
											},
											Name: "result",
										},
									},
								},
								"main": ast.Function{
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
																Name: "identity",
																ArgumentsList: &ast.ArgumentsList{
																	Arguments: []ast.Expression{
																		ast.Literal{
																			VariableType: types.BasicType{
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
					},
				},
			},
		},
		StructFunctions: map[string]types.Function{},
	}
	assert.Equal(t, expectedProgram, program)
}

func TestGenericStruct(t *testing.T) {
	validProgram(t, testcode.GenericStruct)
}

func TestGenericStructInstance(t *testing.T) {
	validProgram(t, testcode.GenericStructInstance1)
	validProgram(t, testcode.GenericStructInstance2)
}

func TestGenericInterfaceFunction(t *testing.T) {
	validProgram(t, testcode.GenericInterfaceFunction)
}

func TestGenericImplementedInterfaceFunctionAllAnnotated(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionAllAnnotated)
}

func TestGenericImplementedInterfaceFunctionAnnotatedReturnType(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionAnnotatedReturnType)
}

func TestGenericImplementedInterfaceFunctionAnnotatedArg(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionAnnotatedArg)
}

func TestGenericImplementedInterfaceFunctionNotAnnotated(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionNotAnnotated)
}
