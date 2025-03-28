package parser_typer_test

import (
	"testing"
)

func TestMainProgramWithWrongArgCount(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime, anotherRuntime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)
`, "expected 1 params but got 2")
}

func TestInvalidMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    applyToString := (f: (String) ~> Void, strF: () ~> String): Void => {
      f(strF())
    }
    output := (): String => {
      "Hello world!"
    }
    applyToString(runtime.console.log, () => {false})
  }
)
`, "expected type String but found Boolean")
}

func TestMainProgramWithVariableWithFunctionWithWrongType(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    output := (): String => {
			
    }
    runtime.console.log(output())
  }
)
`, "empty function block not allowed (maybe you want to return null?)")
}

func TestMainProgramWithArgAnnotatedWrongArg(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: String) => {
    runtime.console.log("Hello world!")
  }
)
`, "in parameter position 0 expected type tenecs.go.Runtime but you have annotated String")
}

func TestMainProgramWithArgAnnotatedWrongReturn(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.go.Main

app := Main(
  main = (runtime): String => {
    runtime.console.log("Hello world!")
  }
)
`, "in return type expected type Void but you have annotated String")
}

func TestNamedArgWrong(t *testing.T) {
	invalidProgram(t, `
package main

f := (a: String, b: String): String => {
  a
}

usage := (): String => {
  f(b = "", "")
}
`, "name of argument should be 'a'")
}
