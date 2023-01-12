package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestMainProgramEmpty(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime) => {
		
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
						Block: []ast.Expression{},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedProgram, program)
}

func TestMainProgramReturningStringInBody(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime) => {
		"can't return string'"
	}
}
`, "expected type Void but found String")
}

func TestMainProgramMultipleMains(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime) => {}
	public main := (runtime) => {}
}
`, "two variables declared in module app with name main")
}

func validProgram(t *testing.T, program string) ast.Program {
	res, err := parser.ParseString(program)
	assert.NoError(t, err)

	p, err := typer.Typecheck(*res)
	assert.NoError(t, err)
	return *p
}

func invalidProgram(t *testing.T, program string, errorMessage string) {
	res, err := parser.ParseString(program)
	if err != nil {
		assert.NoError(t, err)
	}

	_, err = typer.Typecheck(*res)
	assert.Error(t, err, "Didn't get an typererror")
	assert.Equal(t, errorMessage, err.Error())
}
