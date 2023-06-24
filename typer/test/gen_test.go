package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/type_error"
	"testing"
)

//go:generate go run ../test_generate/main.go

func validProgram(t *testing.T, program string) ast.Program {
	res, err := parser.ParseString(program)
	assert.NoError(t, err)

	p, typeErr := typer.Typecheck(*res)
	if typeErr != nil {
		t.Fatal(type_error.Render(program, typeErr.(*type_error.TypecheckError)))
	}
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
