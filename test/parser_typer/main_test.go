package parser_typer_test

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestMainProgramMissingImport(t *testing.T) {
	invalidProgram(`
package main

import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "could not resolve annotated type Runtime")
}

func TestMainProgramWithWrongArgCount(t *testing.T) {
	invalidProgram(`
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime, anotherRuntime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}
`, "expected 1 parameters but got 2")
}

func TestMainProgramWithArgAnnotatedArg(t *testing.T) {
	validProgram(`
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
	invalidProgram(`
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: String) => {
		runtime.console.log("Hello world!")
	}
}
`, "parameter runtime needs to be of type Interface with variables (console) but it's annotated with type String")
}

func TestMainProgramWithArgAnnotatedReturn(t *testing.T) {
	validProgram(`
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
	invalidProgram(`
package main

import tenecs.os.Main

module app: Main {
	public main := (runtime): String => {
		runtime.console.log("Hello world!")
	}
}
`, "Expected lambda return type Void but you annotated String")
}

func TestMainProgramWithArgAnnotatedArgAndReturn(t *testing.T) {
	validProgram(`
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

func validProgram(program string) {
	res, err := parser.ParseString(program)
	if err != nil {
		panic(err)
	}

	err = typer.Validate(*res)
	if err != nil {
		panic(err)
	}
}

func invalidProgram(program string, errorMessage string) {
	res, err := parser.ParseString(program)
	if err != nil {
		panic(err)
	}

	err = typer.Validate(*res)
	if err == nil {
		panic("Didn't get an error")
	}
	if err.Error() != errorMessage {
		panic("Expected error [" + errorMessage + "] but got [" + err.Error() + "]")
	}
}
