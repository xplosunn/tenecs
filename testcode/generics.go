package testcode

const Generics TestCodeCategory = "generics"

var GenericFunctionDeclared = Create(Generics, "GenericFunctionDeclared", `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime): Void => {
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)

var GenericFunctionInvoked1 = Create(Generics, "GenericFunctionInvoked1", `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime): Void => {
		output := "Hello world!"
		hw := identity<String>(output)
		runtime.console.log(hw)
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)

var GenericFunctionInvoked2 = Create(Generics, "GenericFunctionInvoked2", `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime): Void => {
		hw := identity<String>("Hello world!")
		runtime.console.log(hw)
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)

var GenericFunctionInvoked3 = Create(Generics, "GenericFunctionInvoked3", `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime): Void => {
		runtime.console.log(identity<String>("Hello world!"))
	}
	identity := <T>(arg: T): T => {
		arg
	}
}
`)

var GenericFunctionInvoked4 = Create(Generics, "GenericFunctionInvoked4", `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime): Void => {
		output := "Hello world!"
		hw := identity<String>(output)
		runtime.console.log(hw)
	}
	identity := <T>(arg: T): T => {
		result := arg
		result
	}
}
`)

var GenericFunctionDoubleInvoked = Create(Generics, "GenericFunctionDoubleInvoked", `
package main

import tenecs.os.Main

app := (): Main => implement Main {
	public main := (runtime): Void => {
		runtime.console.log(identity<String>("ciao"))
	}
	identity := <T>(arg: T): T => {
		output := identityFn<T>(arg)
		output
	}
	identityFn := <A>(arg: A): A => {
		result := arg
		result
	}
}
`)

var GenericStruct = Create(Generics, "GenericStruct", `
package main

struct Box<T>(inside: T)
`)

var GenericStructInstance = Create(Generics, "GenericStructInstance", `
package main

import tenecs.os.Main

struct Box<T>(inside: T)

app := (): Main => implement Main {
	public main := (runtime) => {
		box := Box<String>("Hello world!")
		runtime.console.log(box.inside)
	}
}
`)

var GenericInterfaceFunction = Create(Generics, "GenericInterfaceFunction", `
package main

interface Assert {
	public equal: <T>(T, T) -> Void
}
`)

var GenericImplementedInterfaceFunctionAllAnnotated = Create(Generics, "GenericImplementedInterfaceFunctionAllAnnotated", `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

id := (): IdentityFunction => implement IdentityFunction {
	public identity := <T>(t: T): T => {
		t
	}
}
`)

var GenericImplementedInterfaceFunctionAnnotatedReturnType = Create(Generics, "GenericImplementedInterfaceFunctionAnnotatedReturnType", `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

id := (): IdentityFunction => implement IdentityFunction {
	public identity := <T>(t): T => {
		t
	}
}
`)

var GenericImplementedInterfaceFunctionAnnotatedArg = Create(Generics, "GenericImplementedInterfaceFunctionAnnotatedArg", `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

id := (): IdentityFunction => implement IdentityFunction {
	public identity := <T>(t: T) => {
		t
	}
}
`)

var GenericImplementedInterfaceFunctionNotAnnotated = Create(Generics, "GenericImplementedInterfaceFunctionNotAnnotated", `
package main

interface IdentityFunction {
	public identity: <T>(T) -> T
}

id := (): IdentityFunction => implement IdentityFunction {
	public identity := <T>(t) => {
		t
	}
}
`)

var GenericFunctionFixingArray = Create(Generics, "GenericFunctionFixingArray", `
package mypackage

import tenecs.array.emptyArray

emptyStringArray := (): Array<String> => {
  emptyArray<String>()
}
`)

var GenericFunctionSingleElementArray = Create(Generics, "GenericFunctionSingleElementArray", `
package mypackage

import tenecs.array.emptyArray
import tenecs.array.append

arrayOf := (elem: String): Array<String> => {
  append<String>(emptyArray<String>(), elem)
}
`)

var GenericFunctionTakingArray = Create(Generics, "GenericFunctionSingleElementArray", `
package mypackage

toJson := <T>(t: T): String => {
  "not actually implemented"
}

doStuff := (): String => {
  arr := [String]("a", "b")
  toJson<Array<String>>(arr)
}
`)

var GenericStructFunction = Create(Generics, "GenericStructFunction", `
package mypackage

struct Box<T>(elem: T)

f := <T>(): Box<String> => {
  b := Box<String>("wee")
  b
}
`)
