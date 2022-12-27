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
