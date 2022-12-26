package parser_typer_test

import "testing"

func TestMainProgramWithVariable(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		output := "Hello world!"
		runtime.console.log(output)
	}
}
`)
}
