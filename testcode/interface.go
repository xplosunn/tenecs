package testcode

const Interface TestCodeCategory = "interface"

var InterfaceEmpty = Create(Interface, "InterfaceEmpty", `
package main

interface A {
}
`)

var InterfaceWithSeparateModuleEmpty1 = Create(Interface, "InterfaceWithSeparateModuleEmpty1", `
package main

interface A {
}

a := (): A => implement A {
}
`)

var InterfaceWithSeparateModuleEmpty2 = Create(Interface, "InterfaceWithSeparateModuleEmpty2", `
package main

a := (): A => implement A {
}

interface A {
}
`)

var InterfaceVariableString = Create(Interface, "InterfaceVariableString", `
package main

interface A {
	public a: String
}
`)

var InterfaceVariableFunctionZeroArgs = Create(Interface, "InterfaceVariableFunctionZeroArgs", `
package main

interface A {
	public a: () -> String
}
`)

var InterfaceVariableFunctionOneArg = Create(Interface, "InterfaceVariableFunctionOneArg", `
package main

interface A {
	public a: (String) -> String
}
`)

var InterfaceVariableFunctionTwoArgs = Create(Interface, "InterfaceVariableFunctionTwoArgs", `
package main

interface A {
	public a: (String, Boolean) -> String
}
`)

var InterfaceWithSeparateModuleVariableString = Create(Interface, "InterfaceWithSeparateModuleVariableString", `
package main

interface A {
	public a: String
}

app := (): A => implement A {
	public a := ""
}
`)

var InterfaceReturningAnotherInterfaceInVariable = Create(Interface, "InterfaceReturningAnotherInterfaceInVariable", `
package main

interface Goods {
	public name: String
}

interface Factory {
	public produce: () -> Goods
}
`)
