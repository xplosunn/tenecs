package parser_typer_test

import "testing"

func TestMainProgramMissingBothImports(t *testing.T) {
	invalidProgram(t, `
package main

app := (): Main => implement Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "not found type: Main")
}

func TestMainProgramMissingOneImport(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "not found type: Runtime")
}
