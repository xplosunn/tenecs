package parser_typer_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

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

func TestMainProgramWithWrongArgCount(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime, anotherRuntime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "expected same number of arguments as interface variable (1) but found 2")
}

func TestMainProgramWithArgAnnotatedArg(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithArgAnnotatedWrongArg(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: String) => {
		runtime.console.log("Hello world!")
	}
}
`, "in parameter position 0 expected type tenecs.os.Runtime but you have annotated String")
}

func TestMainProgramWithArgAnnotatedReturn(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main

module app: Main {
	public main := (runtime): Void => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithArgAnnotatedWrongReturn(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Main

module app: Main {
	public main := (runtime): String => {
		runtime.console.log("Hello world!")
	}
}
`, "in return type expected type Void but you have annotated String")
}

func TestMainProgramWithArgAnnotatedArgAndReturn(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime): Void => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithAnotherFunctionTakingRuntime(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main
import tenecs.os.Runtime

module app: Main {
	public main := (runtime) => {
		mainRun(runtime)
	}
	mainRun := (runtime: Runtime): Void => {
		runtime.console.log("Hello world!")
	}
}
`)
}

func TestMainProgramWithAnotherFunctionTakingConsole(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Main
import tenecs.os.Runtime
import tenecs.os.Console

module app: Main {
	public main := (runtime) => {
		mainRun(runtime.console)
	}
	mainRun := (console: Console): Void => {
		console.log("Hello world!")
	}
}
`)
}

func validProgram(t *testing.T, program string) {
	res, err := parser.ParseString(program)
	assert.NoError(t, err)

	err = typer.Typecheck(*res)
	assert.NoError(t, err)
}

func invalidProgram(t *testing.T, program string, errorMessage string) {
	res, err := parser.ParseString(program)
	if err != nil {
		assert.NoError(t, err)
	}

	err = typer.Typecheck(*res)
	assert.Error(t, err, "Didn't get an error")
	assert.Equal(t, errorMessage, err.Error())
}
