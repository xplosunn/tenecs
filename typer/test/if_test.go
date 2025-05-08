package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/typer/ast"
	"testing"
)

func TestMainProgramWithIfNonBooleanCondition(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    if runtime {
      runtime.console.log("Hello world!")
    }
  }
)
`, "expected type Boolean but found tenecs.go.Runtime")
}

func TestIfElse(t *testing.T) {
	ast1 := validProgram(t, `
package pkg

f := (cond: (String) ~> Boolean): String => {
  if cond("a") {
    if cond("a1") {
      null
    }
    "x"
  } else if cond("b") {
    if cond("b1") {
      null
    } else if cond("b2") {
      null
    }
    "y"
  } else {
    "z"
  }
}
`)
	ast2 := validProgram(t, `
package pkg

f := (cond: (String) ~> Boolean): String => {
  if cond("a") {
    if cond("a1") {
      null
    }
    "x"
  } else {
    if cond("b") {
      if cond("b1") {
        null
      } else {
        if cond("b2") {
          null
        }
      }
      "y"
    } else {
      "z"
    }
  }
}
`)
	assert.Equal(t, ast.EmptyCodePoints(ast1), ast.EmptyCodePoints(ast2))
}
