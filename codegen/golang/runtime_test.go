package golang_test

import (
	"github.com/alecthomas/assert/v2"
	golang2 "github.com/xplosunn/tenecs/codegen/golang"
	"github.com/xplosunn/tenecs/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
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

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := golang2.GenerateProgramMain(typed, nil)

	output := golang.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}
