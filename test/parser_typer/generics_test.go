package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestGenericFunctionDeclared(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): Void => {
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)
}

func TestGenericFunctionInvoked(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): Void => {
		output := "Hello world!"
		hw := identity<String>(output)
		runtime.console.log(hw)
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)
	validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): Void => {
		hw := identity<String>("Hello world!")
		runtime.console.log(hw)
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)
	validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): Void => {
		runtime.console.log(identity<String>("Hello world!"))
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)
	program := validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): Void => {
		output := "Hello world!"
		hw := identity<String>(output)
		runtime.console.log(hw)
	}
	identity := <T>(arg: T): T => {
		result := arg
		result
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
					"identity": &ast.Function{
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
								Expression: ast.ReferenceOrInvocation{
									VariableType: types.TypeArgument{
										Name: "T",
									},
									DotSeparatedVars: []string{
										"arg",
									},
								},
							},
							ast.ReferenceOrInvocation{
								VariableType: types.TypeArgument{
									Name: "T",
								},
								DotSeparatedVars: []string{
									"result",
								},
							},
						},
					},
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
							ast.Declaration{
								VariableType: types.Void{},
								Name:         "hw",
								Expression: ast.ReferenceOrInvocation{
									VariableType: types.BasicType{
										Type: "String",
									},
									DotSeparatedVars: []string{
										"identity",
									},
									Arguments: &ast.ArgumentsList{
										Arguments: []ast.Expression{
											ast.ReferenceOrInvocation{
												VariableType: types.BasicType{
													Type: "String",
												},
												DotSeparatedVars: []string{
													"output",
												},
											},
										},
									},
								},
							},
							ast.ReferenceOrInvocation{
								VariableType: types.Void{},
								DotSeparatedVars: []string{
									"runtime",
									"console",
									"log",
								},
								Arguments: &ast.ArgumentsList{
									Arguments: []ast.Expression{
										ast.ReferenceOrInvocation{
											VariableType: types.BasicType{
												Type: "String",
											},
											DotSeparatedVars: []string{
												"hw",
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

func TestGenericFunctionDoubleInvoked(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime): Void => {
		runtime.console.log(identity<String>("ciao"))
	}
	identity := <T>(arg: T): T => {
		output := identityFn<T>(arg)
		output
	}
	identityFn := <A>(arg: A): A => {
		result := arg
		result
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
					"identity": &ast.Function{
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
								Expression: ast.ReferenceOrInvocation{
									VariableType: types.TypeArgument{
										Name: "T",
									},
									DotSeparatedVars: []string{
										"identityFn",
									},
									Arguments: &ast.ArgumentsList{
										Arguments: []ast.Expression{
											ast.ReferenceOrInvocation{
												VariableType: types.TypeArgument{
													Name: "T",
												},
												DotSeparatedVars: []string{
													"arg",
												},
											},
										},
									},
								},
							},
							ast.ReferenceOrInvocation{
								VariableType: types.TypeArgument{
									Name: "T",
								},
								DotSeparatedVars: []string{
									"output",
								},
							},
						},
					},
					"identityFn": &ast.Function{
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
								Expression: ast.ReferenceOrInvocation{
									VariableType: types.TypeArgument{
										Name: "A",
									},
									DotSeparatedVars: []string{
										"arg",
									},
								},
							},
							ast.ReferenceOrInvocation{
								VariableType: types.TypeArgument{
									Name: "A",
								},
								DotSeparatedVars: []string{
									"result",
								},
							},
						},
					},
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
							ast.ReferenceOrInvocation{
								VariableType: types.Void{},
								DotSeparatedVars: []string{
									"runtime",
									"console",
									"log",
								},
								Arguments: &ast.ArgumentsList{
									Arguments: []ast.Expression{
										ast.ReferenceOrInvocation{
											VariableType: types.BasicType{
												Type: "String",
											},
											DotSeparatedVars: []string{
												"identity",
											},
											Arguments: &ast.ArgumentsList{
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
	}
	assert.Equal(t, expectedProgram, program)
}

func TestGenericStruct(t *testing.T) {
	validProgram(t, `
package main

struct Box<T>(inside: T)
`)
}

func TestGenericStructInstance(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main

struct Box<T>(inside: T)

implementing Main module app {
	public main := (runtime) => {
		box := Box<String>("Hello world!")
		runtime.console.log(box.inside)
	}
}
`)
	validProgram(t, `
package main

import tenecs.os.Main

struct Box<T>(inside: T)

implementing Main module app {
	public main := (runtime) => {
		box := Box<String>("Hello world!")
		runtime.console.log(box.inside)
	}
}
`)
}

func TestGenericInterfaceFunction(t *testing.T) {
	validProgram(t, `
package main

interface Assert {
	public equal: <T>(T, T) -> Void
}
`)
}

func TestGenericImplementedInterfaceFunctionAllAnnotated(t *testing.T) {
	validProgram(t, `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

implementing IdentityFunction module id {
	public identity := <T>(t: T): T => {
		t
	}
}
`)
}

func TestGenericImplementedInterfaceFunctionAnnotatedReturnType(t *testing.T) {
	validProgram(t, `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

implementing IdentityFunction module id {
	public identity := <T>(t): T => {
		t
	}
}
`)
}

func TestGenericImplementedInterfaceFunctionAnnotatedArg(t *testing.T) {
	validProgram(t, `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

implementing IdentityFunction module id {
	public identity := <T>(t: T) => {
		t
	}
}
`)
}

func TestGenericImplementedInterfaceFunctionNotAnnotated(t *testing.T) {
	validProgram(t, `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

implementing IdentityFunction module id {
	public identity := <T>(t) => {
		t
	}
}
`)
}
