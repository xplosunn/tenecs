package ast_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"testing"
)

func TestDetermineRefDependencies(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main
import tenecs.int.times
import tenecs.int.minus
import tenecs.compare.eq
import tenecs.json.jsonInt

factorial := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    prev := factorial(minus(i, 1))
    times(i, prev)
  }
}

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log(jsonInt().toJson(factorial(5)))
  }
)`
	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	desugared := desugar.Desugar(*parsed)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)

	refDependencies := ast.DetermineRefDependencies(*typed)

	expectedFactorialDependencies := ast.Set[ast.Ref]{}
	expectedFactorialDependencies.Put(ast.Ref{
		Package: "tenecs.compare",
		Name:    "eq",
	})
	expectedFactorialDependencies.Put(ast.Ref{
		Package: "tenecs.int",
		Name:    "minus",
	})
	expectedFactorialDependencies.Put(ast.Ref{
		Package: "tenecs.int",
		Name:    "times",
	})

	assert.Equal(t, expectedFactorialDependencies, refDependencies[ast.Ref{
		Package: "main",
		Name:    "factorial",
	}])

	expectedAppDependencies := ast.Set[ast.Ref]{}
	expectedAppDependencies.Put(ast.Ref{
		Package: "main",
		Name:    "factorial",
	})
	expectedAppDependencies.Put(ast.Ref{
		Package: "tenecs.go",
		Name:    "Console",
	})
	expectedAppDependencies.Put(ast.Ref{
		Package: "tenecs.go",
		Name:    "Main",
	})
	expectedAppDependencies.Put(ast.Ref{
		Package: "tenecs.go",
		Name:    "Runtime",
	})
	expectedAppDependencies.Put(ast.Ref{
		Package: "tenecs.json",
		Name:    "JsonConverter",
	})
	expectedAppDependencies.Put(ast.Ref{
		Package: "tenecs.json",
		Name:    "jsonInt",
	})

	assert.Equal(t, expectedAppDependencies, refDependencies[ast.Ref{
		Package: "main",
		Name:    "app",
	}])
}
