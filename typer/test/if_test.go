package parser_typer_test

import (
	"testing"
)

func TestMainProgramWithIfNonBooleanCondition(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime: Runtime) => {
		if runtime {
			runtime.console.log("Hello world!")
		}
	}
}
`, "expected type Boolean but found tenecs.os.Runtime")
}
