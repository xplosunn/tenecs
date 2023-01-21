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

implementing A module app(arg: NonExistingType) {
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

func TestModuleWithConstructorWithArgPublic(t *testing.T) {
	validProgram(t, testcode.ModuleWithConstructorWithArgPublic)
}

func TestModuleWithConstructorWithSameNameAsArg(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

implementing A module a(public a: String) {
	
}
`, "variable a cannot have the same name as the module")
}

func TestModuleWithConstructorWithSameNameAsVariable(t *testing.T) {
	invalidProgram(t, `
package main

interface A {
	public a: String
}

implementing A module a() {
	public a := ""
}
`, "variable a cannot have the same name as the module")
}

func TestModuleCreation(t *testing.T) {
	validProgram(t, testcode.ModuleCreation1)
	validProgram(t, testcode.ModuleCreation2)
	validProgram(t, testcode.ModuleCreation3)
}

func TestModuleSelfCreation(t *testing.T) {
	validProgram(t, testcode.ModuleSelfCreation)
}
