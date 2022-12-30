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

implementing A module a {
}
`)

	validProgram(t, `
package main

implementing A module a {
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

implementing A module app {
	public a := ""
}
`)
}

func TestInterfaceWithSeparateModuleVariableStringThatShouldBePublic(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

implementing A module app {
	a := ""
}
`, "variable a should be public as it's in implemented interface A")
}

func TestInterfaceWithSeparateModuleVariableStringThatShouldNotBePublic(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

implementing A module app {
	public a := ""
	public b := ""
}
`, "variable b can't be public as no implemented interface has a variable with the same name")
}

func TestInterfaceWithSeparateModuleMissingVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

implementing A module a {
	
}
`, "variable a of interface A missing in module a")
}

func TestInterfaceWithSeparateModuleWrongVariableType(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: Void
}

implementing A module app {
	public a := ""
}
`, "expected type Void but found String")
}
