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
						ast.Implementation{
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
										ReturnType: types.Void(),
									},
									Block: []ast.Expression{
										ast.Declaration{
											Name: "output",
											Expression: ast.Literal{
												VariableType: types.String(),
												Literal: parser.LiteralString{
													Value: "\"Hello world!\"",
												},
											},
										},
										ast.Declaration{
											Name: "hw",
											Expression: ast.Invocation{
												VariableType: types.String(),
												Over: ast.Reference{
													VariableType: &types.Function{
														Arguments: []types.FunctionArgument{
															{
																Name:         "arg",
																VariableType: types.String(),
															},
														},
														ReturnType: types.String(),
													},
													Name: "identity",
												},
												Generics: []types.VariableType{
													types.String(),
												},
												Arguments: []ast.Expression{
													ast.Reference{
														VariableType: types.String(),
														Name:         "output",
													},
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
													VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Console"),
													Over: ast.Reference{
														VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
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
													Name:         "hw",
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
	program.FieldsByType = nil
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
						ast.Implementation{
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
												Generics: []types.VariableType{
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
										ReturnType: types.Void(),
									},
									Block: []ast.Expression{
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
													VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Console"),
													Over: ast.Reference{
														VariableType: standard_library.StdLibGetOrPanic(t, "tenecs.os.Runtime"),
														Name:         "runtime",
													},
													Access: "console",
												},
												Access: "log",
											},
											Generics: []types.VariableType{},
											Arguments: []ast.Expression{
												ast.Invocation{
													VariableType: types.String(),
													Over: ast.Reference{
														VariableType: &types.Function{
															Arguments: []types.FunctionArgument{
																{
																	Name:         "arg",
																	VariableType: types.String(),
																},
															},
															ReturnType: types.String(),
														},
														Name: "identity",
													},
													Generics: []types.VariableType{
														types.String(),
													},
													Arguments: []ast.Expression{
														ast.Literal{
															VariableType: types.String(),
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
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}

func TestWrongGeneric(t *testing.T) {
	invalidProgram(t, `
package mypackage

struct Tuple<L, R>(left: L, right: R)

leftAs := <L, R, T>(tuple: Tuple<L, R>, as: T): Tuple<T, R> => {
  result := Tuple<T, T>(as, as)
  result
}
`, "expected type mypackage.Tuple<<T>, <R>> but found mypackage.Tuple<<T>, <T>>")
}

func TestWrongGeneric2(t *testing.T) {
	invalidProgram(t, `
package mypackage

struct Tuple<L, R>(left: L, right: R)

leftAs := <L, R, T>(tuple: Tuple<L, R>, as: T): Tuple<T, R> => {
  Tuple<T, T>(as, as)
}
`, "expected type mypackage.Tuple<<T>, <R>> but found mypackage.Tuple<<T>, <T>>")
}

func TestGenericFunctionInvocation(t *testing.T) {
	validProgram(t, `
package mypackage

takeArray := <A>(arr: Array<A>): Void => {}

usage := (): Void => {
  takeArray<String | Int>([Int | String]("", 1))
  null
}
`)
}

func TestGenericFunctionInvocation2(t *testing.T) {
	validProgram(t, `
package mypackage

take := <A>(a: A): Void => {}

usage := (): Void => {
  take<Array<String> | String>([String]())
  null
}
`)
}

func TestGenericFunctionInvocation3(t *testing.T) {
	validProgram(t, `
package mypackage

interface Parser<T> {}

parseArray := <Of>(parserOf: Parser<Of>): Parser<Array<Of>> => {
  implement Parser<Array<Of>> {
  }
}

parseString := (): Parser<String> => {
  implement Parser<Array<String>> {
  }
}

takeParser := <Of>(parser: Parser<Of>): Void => {}

usage := (): Void => {
  takeParser<Array<Array<String>>>(parseArray<Array<String>>(parseArray<String>(parseString())))
}

`)
}

func TestGenericFunctionInvocation4(t *testing.T) {
	validProgram(t, `
package mypackage

wrapFunction := <R>(f: () -> R): () -> R => {
  (): R => {
    f()
  }
}

usage := (): Void => {
  f := wrapFunction<Void>(() => null)
  f()
}
`)
}

func TestGenericFunctionInvocation5(t *testing.T) {
	validProgram(t, `
package mypackage

apply := <A, B>(a: A, f: (A) -> B): B => {
  f(a)
}

usage := (): String => {
  apply(1, (int: Int): String => {""})
}
`)
}

func TestGenericFunctionWrongInvocation(t *testing.T) {
	invalidProgram(t, `
package mypackage

take := <A>(arg: A): Void => {}

usage := (): Void => {
  take<String>(1)
  null
}

`, "expected type String but found Int")
}

func TestGenericFunctionWrongInvocation2(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeArray := <A>(arr: Array<A>): Void => {}

usage := (): Void => {
  takeArray<String>([Int](1))
  null
}

`, "expected Array<String> but got Array<Int>")
}

func TestGenericFunctionWrongInvocation3(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeArray := <A>(arr: Array<A>): Void => {}

usage := (): Void => {
  takeArray<String>([String | Int](""))
  null
}

`, "expected Array<String> but got Array<String | Int>")
}

func TestGenericFunctionWrongInvocation4(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeArray := <A>(arr: Array<A>): Void => {}

usage := (): Void => {
  takeArray<Array<String>>([String]())
  null
}

`, "expected Array<Array<String>> but got Array<String>")
}

func TestGenericFunctionWrongInvocation5(t *testing.T) {
	invalidProgram(t, `
package mypackage

assertEqual := <T> (a: T, b: T): Void => {}

arrayOfStringOrString := (): Array<String> | String => {
  ""
}

usage := (): Void => {
	assertEqual<Array<String>>([String](), arrayOfStringOrString())
}

`, "expected type Array<String> but found Array<String> | String")
}
