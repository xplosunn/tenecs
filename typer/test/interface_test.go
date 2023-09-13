package parser_typer_test

import (
	"testing"
)

func TestInterfaceVariablesSameName(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
	public a: String
}
`, "more than one variable with name 'a'")
}

func TestInterfaceWithSeparateImplementationVariableStringThatShouldBePublic(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

app := (): A => implement A {
	a := ""
}
`, "variable a should be public")
}

func TestInterfaceWithSeparateImplementationVariableStringThatShouldNotBePublic(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

app := (): A => implement A {
	public a := ""
	public b := ""
}
`, "variable b should not be public")
}

func TestInterfaceWithSeparateImplementationMissingVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

app := ():A => implement A {
	
}
`, "missing declaration for variable a")
}

func TestInterfaceWithSeparateImplementationWrongVariableType(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: Void
}

app := (): A => implement A {
	public a := ""
}
`, "expected type Void but found String")
}
