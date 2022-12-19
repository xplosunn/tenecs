package parser_typer_test

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestMainProgramWithArgAnnotatedArg(t *testing.T) {
	validProgram(`
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func validProgram(program string) {
	res, err := parser.ParseString(program)
	if err != nil {
		panic(err)
	}

	err = typer.Validate(*res)
	if err != nil {
		panic(err)
	}
}
