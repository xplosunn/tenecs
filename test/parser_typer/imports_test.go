package parser_typer_test

import "testing"

func TestMainProgramMissingBothImports(t *testing.T) {
	invalidProgram(t, `
package main

implementing Main module app {
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

implementing Main module app {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "not found type: Runtime")
}
