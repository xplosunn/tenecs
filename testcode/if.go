package testcode

const If TestCodeCategory = "if"

var MainProgramWithIf = Create(If, "MainProgramWithIf", `
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app {
	public main := (runtime: Runtime) => {
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

implementing Main module app {
	public main := (runtime: Runtime) => {
		if false {
			runtime.console.log("Hello world!")
		} else {
			runtime.console.log("Hello world!")
		}
	}
}
`)
