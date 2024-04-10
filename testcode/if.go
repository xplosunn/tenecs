package testcode

const If TestCodeCategory = "if"

var MainProgramWithIf = Create(If, "MainProgramWithIf", `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
	main := (runtime: Runtime) => {
		if true {
			runtime.console.log("Hello world!")
		}
	}
}
`)

var MainProgramWithIfElse = Create(If, "MainProgramWithIfElse", `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
  main := (runtime: Runtime) => {
		if false {
			runtime.console.log("Hello world!")
		} else {
			runtime.console.log("Hello world!")
		}
	}
}
`)

var MainProgramWithIfElseIf = Create(If, "MainProgramWithIfElseIf", `
package main

import tenecs.os.Runtime
import tenecs.os.Main

app := (): Main => implement Main {
  main := (runtime: Runtime) => {
		if false {
			runtime.console.log("Hello world!")
		} else if false {
			runtime.console.log("Hello world!")
		} else if true {
			runtime.console.log("Hello world!")
		} else {
			runtime.console.log("Hello world!")
		}
	}
}
`)
