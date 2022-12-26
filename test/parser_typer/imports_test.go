package parser_typer_test

import "testing"

func TestMainProgramMissingBothImports(t *testing.T) {
	invalidProgram(t, `
package main

module app: Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "not found interface with name Main")
}

func TestMainProgramMissingOneImport(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "not found type: Runtime")
}
