package testcode

const Interface TestCodeCategory = "interface"

var InterfaceEmpty = Create(Interface, "InterfaceEmpty", `
package main

interface A {
}
`)

var InterfaceWithSeparateImplementationEmpty1 = Create(Interface, "InterfaceWithSeparateImplementationEmpty1", `
package main

interface A {
}

a := (): A => implement A {
}
`)

var InterfaceWithSeparateImplementationEmpty2 = Create(Interface, "InterfaceWithSeparateImplementationEmpty2", `
package main

a := (): A => implement A {
}

interface A {
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

var InterfaceReturningAnotherInterfaceInVariable = Create(Interface, "InterfaceReturningAnotherInterfaceInVariable", `
package main

interface Goods {
	public name: () -> String
}

interface Factory {
	public produce: () -> Goods
}
`)
