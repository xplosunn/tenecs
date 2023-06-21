package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestModuleWithConstructorEmpty(t *testing.T) {
	validProgram(t, testcode.ModuleWithConstructorEmpty)
}

func TestModuleWithInvalidType(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

app := (arg: NonExistingType): A => implement A {
	public a := ""
}
`, "not found type: NonExistingType")
}

func TestModuleWithConstructorWithArgUnused(t *testing.T) {
	validProgram(t, testcode.ModuleWithConstructorWithArgUnused)
}

func TestModuleWithConstructorWithArgUsed(t *testing.T) {
	validProgram(t, testcode.ModuleWithConstructorWithArgUsed)
}

func TestModuleWithConstructorWithSameNameAsVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

a := (): A => implement A {
	public a := ""
}
`, "duplicate variable 'a'")
}

func TestModuleCreation(t *testing.T) {
	validProgram(t, testcode.ModuleCreation1)
	validProgram(t, testcode.ModuleCreation2)
	validProgram(t, testcode.ModuleCreation3)
}

func TestModuleSelfCreation(t *testing.T) {
	validProgram(t, testcode.ModuleSelfCreation)
}
