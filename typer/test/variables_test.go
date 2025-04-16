package parser_typer_test

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func TestMainProgramWithVariable(t *testing.T) {
	program := validProgram(t, `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    output := "Hello world!"
    runtime.console.log(output)
  }
)
`)
	snaps.MatchStandaloneSnapshot(t, program)
}

func TestInvalidVariableName(t *testing.T) {
	program := `package pk

true := false
`

	invalidProgram(t, program, "Variable can't be named 'true'")
}

func TestInvalidLocalVariableName(t *testing.T) {
	program := `package pk

_ := (): Void => {
  false := true
}
`

	invalidProgram(t, program, "Variable can't be named 'false'")
}
