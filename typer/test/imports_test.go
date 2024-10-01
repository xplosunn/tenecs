package parser_typer_test

import "testing"

func TestMainProgramMissingBothImports(t *testing.T) {
	invalidProgram(t, `
package main

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)
`, "Reference not found: Main")
}

func TestMainProgramMissingOneImport(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)
`, "not found type: Runtime")
}
