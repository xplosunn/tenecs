package codegen_golang_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"testing"
)

func TestRef(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    ref := runtime.ref.new("1st value")
    runtime.console.log(ref.get())
    ref.set("2nd value")
    runtime.console.log(ref.get())
  }
)`
	expectedRunResult := `1st value
2nd value
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared := desugar.Desugar(*parsed)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, ast.Ref{
		Package: "main",
		Name:    "app",
	})

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}
