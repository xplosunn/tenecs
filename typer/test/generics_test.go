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
	mainStr := "main"
	expectedProgram := ast.Program{
		Declarations: map[ast.Ref]ast.Expression{
			ast.Ref{
				Package: "main",
				Name:    "app",
			}: mainWithBlock(t, []ast.Expression{
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
							PackageName: &mainStr,
							Name:        "identity",
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
							Name:         "hw",
						},
					},
				},
			}),
			ast.Ref{
				Package: "main",
				Name:    "identity",
			}: &ast.Function{
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
		},
		StructFunctions: map[ast.Ref]*types.Function{},
		NativeFunctions: map[ast.Ref]*types.Function{
			ast.Ref{
				Package: "tenecs_go",
				Name:    "Main",
			}: mainNativeFunction(),
		},
	}
	program.FieldsByType = nil
	assert.Equal(t, expectedProgram, program)
}

func TestExpectedGenericFunctionDoubleInvoked(t *testing.T) {
	program := validProgram(t, testcode.GenericFunctionDoubleInvoked)
	mainStr := "main"
	expectedProgram := ast.Program{
		Declarations: map[ast.Ref]ast.Expression{
			ast.Ref{
				Package: "main",
				Name:    "app",
			}: mainWithBlock(t, []ast.Expression{
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
								PackageName: &mainStr,
								Name:        "identity",
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
			}),
			ast.Ref{
				Package: "main",
				Name:    "identity",
			}: &ast.Function{
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
								PackageName: &mainStr,
								Name:        "identityFn",
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
			ast.Ref{
				Package: "main",
				Name:    "identityFn",
			}: &ast.Function{
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
		},
		StructFunctions: map[ast.Ref]*types.Function{},
		NativeFunctions: map[ast.Ref]*types.Function{
			ast.Ref{
				Package: "tenecs_go",
				Name:    "Main",
			}: mainNativeFunction(),
		},
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

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<String | Int>(<Int | String>["", 1])
  null
}
`)
}

func TestGenericFunctionInvocation2(t *testing.T) {
	validProgram(t, `
package mypackage

take := <A>(a: A): Void => {
  null
}

usage := (): Void => {
  take<List<String> | String>(<String>[])
  null
}
`)
}

func TestGenericFunctionInvocation3(t *testing.T) {
	validProgram(t, `
package mypackage

struct Parser<T>()

parseList := <Of>(parserOf: Parser<Of>): Parser<List<Of>> => {
  Parser<List<Of>>()
}

parseString := (): Parser<String> => {
  Parser<String>()
}

takeParser := <Of>(parser: Parser<Of>): Void => {
  null
}

usage := (): Void => {
  takeParser<List<List<String>>>(parseList<List<String>>(parseList<String>(parseString())))
}

`)
}

func TestGenericFunctionInvocation4(t *testing.T) {
	validProgram(t, `
package mypackage

wrapFunction := <R>(f: () ~> R): () ~> R => {
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

apply := <A, B>(a: A, f: (A) ~> B): B => {
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

take := <A>(arg: A): Void => {
  null
}

usage := (): Void => {
  take<String>(1)
  null
}

`, "expected type String but found Int")
}

func TestGenericFunctionWrongInvocation2(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<String>(<Int>[1])
  null
}

`, "expected List<String> but got List<Int>")
}

func TestGenericFunctionWrongInvocation3(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<String>(<String | Int>[""])
  null
}

`, "expected List<String> but got List<String | Int>")
}

func TestGenericFunctionWrongInvocation4(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<List<String>>(<String>[])
  null
}

`, "expected List<List<String>> but got List<String>")
}

func TestGenericFunctionWrongInvocation5(t *testing.T) {
	invalidProgram(t, `
package mypackage

assertEqual := <T> (a: T, b: T): Void => {
  null
}

listOfStringOrString := (): List<String> | String => {
  ""
}

usage := (): Void => {
	assertEqual<List<String>>(<String>[], listOfStringOrString())
}

`, "expected type List<String> but found List<String> | String")
}
