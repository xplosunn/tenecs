package codegen_golang_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/xplosunn/tenecs/codegen2/codegen_golang"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestGenerateProgramNonRunnableMain(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)`
	expectedGoCode := `package main

import ()

func main__app() any {
	return tenecs_go__Main().(func(any) any)(func(_runtime any) any {
		return _runtime.(map[string]any)["_console"].(map[string]any)["_log"].(func(any) any)("Hello world!")
	})
}

func tenecs_go__Main() any {
	log := func(msg any) any {
		println(msg.(string))
		return nil
	}
	console := map[string]any{
		"_log": log,
	}
	runtime := map[string]any{
		"_console": console,
	}
	return func(run any) any {
		return run.(func(any) any)(runtime)
	}
}

func main() {
	main__app()
}
`

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared := desugar.Desugar(*parsed)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	codeIR := ir.ToIR(*typed)

	mainPackage := "main"
	generated := codegen_golang.GenerateProgramMain(&codeIR, ir.Reference{
		Name: ir.VariableName(&mainPackage, "app"),
	})
	generatedFormatted := golang.Fmt(t, generated)
	assert.Equal(t, expectedGoCode, generatedFormatted)

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)

	snaps.MatchStandaloneSnapshot(t, golang.Fmt(t, generated))
}
