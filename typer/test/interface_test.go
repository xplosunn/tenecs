package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestInterfaceEmpty(t *testing.T) {
	validProgram(t, testcode.InterfaceEmpty)
}

func TestInterfaceWithSeparateModuleEmpty(t *testing.T) {
	validProgram(t, testcode.InterfaceWithSeparateModuleEmpty1)
	validProgram(t, testcode.InterfaceWithSeparateModuleEmpty2)
}

func TestInterfaceVariableString(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableString)
}

func TestInterfaceVariableFunctionZeroArgs(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableFunctionZeroArgs)
}

func TestInterfaceVariableFunctionOneArg(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableFunctionOneArg)
}

func TestInterfaceVariableFunctionTwoArgs(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableFunctionTwoArgs)
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
	validProgram(t, testcode.InterfaceWithSeparateModuleVariableString)
}

func TestInterfaceWithSeparateModuleVariableStringThatShouldBePublic(t *testing.T) {
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

func TestInterfaceWithSeparateModuleVariableStringThatShouldNotBePublic(t *testing.T) {
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

func TestInterfaceWithSeparateModuleMissingVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

app := ():A => implement A {
	
}
`, "missing declaration for variable a")
}

func TestInterfaceWithSeparateModuleWrongVariableType(t *testing.T) {
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

func TestInterfaceReturningAnotherInterfaceInVariable(t *testing.T) {
	validProgram(t, testcode.InterfaceReturningAnotherInterfaceInVariable)
}
