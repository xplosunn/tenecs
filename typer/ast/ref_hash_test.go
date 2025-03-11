package ast_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"testing"
)

func TestDetermineRefHashes(t *testing.T) {
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

factorialSameImpl := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    prev := factorial(minus(i, 1))
    times(i, prev)
  }
}

factorialOtherImpl := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    times(i, factorial(minus(i, 1)))
  }
}

`
	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	refHashes, err := ast.DetermineRefHashes(*typed)
	assert.NoError(t, err)
	assert.Equal(t, "WwJ-VUVn3e0Nb3JsonL9RRFwaLA=", refHashes[ast.Ref{
		Package: "main",
		Name:    "factorial",
	}])
	assert.Equal(t, "WwJ-VUVn3e0Nb3JsonL9RRFwaLA=", refHashes[ast.Ref{
		Package: "main",
		Name:    "factorialSameImpl",
	}])
	assert.Equal(t, "QEXF6nGf4Px5WDy4YCp3P5-kUbg=", refHashes[ast.Ref{
		Package: "main",
		Name:    "factorialOtherImpl",
	}])
}
