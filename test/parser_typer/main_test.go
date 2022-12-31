package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestMainProgramEmpty(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main

implementing Main module app {
	public main := (runtime) => {
		
	}
}
`)
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

func validProgram(t *testing.T, program string) {
	res, err := parser.ParseString(program)
	assert.NoError(t, err)

	_, err = typer.Typecheck(*res)
	assert.NoError(t, err)
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
