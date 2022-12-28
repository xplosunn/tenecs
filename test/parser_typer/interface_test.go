package parser_typer_test

import "testing"

func TestInterfaceEmpty(t *testing.T) {
	validProgram(t, `
package main

interface A {
}
`)
}

func TestInterfaceWithSeparateModuleEmpty(t *testing.T) {
	validProgram(t, `
package main

interface A {
}

module a: A {
}
`)

	validProgram(t, `
package main

module a: A {
}

interface A {
}
`)
}

func TestInterfaceVariableString(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}
`)
}

func TestInterfaceVariableFunctionZeroArgs(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: () -> String
}
`)
}

func TestInterfaceVariableFunctionOneArg(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: (String) -> String
}
`)
}

func TestInterfaceVariableFunctionTwoArgs(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: (String, Boolean) -> String
}
`)
}

func TestInterfaceVariablesSameName(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
	public a: String
}
`, "more than one variable with name 'a'")
}

func TestInterfaceWithSeparateModuleVariableString(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}

module a: A {
	a := ""
}
`)
}

func TestInterfaceWithSeparateModuleWrongVariableType(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: Void
}

module a: A {
	public a := ""
}
`, "expected type Void but found String")
}
