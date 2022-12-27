package parser_typer_test

import "testing"

func TestMainProgramWithIf(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		if true {
			runtime.console.log("Hello world!")
		}
	}
}
`)
}

func TestMainProgramWithIfNonBooleanCondition(t *testing.T) {
	invalidProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		if runtime {
			runtime.console.log("Hello world!")
		}
	}
}
`, "in expression 'runtime' expected Boolean but found tenecs.os.Runtime")
}

func TestMainProgramWithIfElse(t *testing.T) {
	validProgram(t, `
package main

import tenecs.os.Runtime
import tenecs.os.Main

module app: Main {
	public main := (runtime: Runtime) => {
		if false {
			runtime.console.log("Hello world!")
		} else {
			runtime.console.log("Hello world!")
		}
	}
}
`)
}
