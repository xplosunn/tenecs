package parser_typer_test

import "testing"

func TestModuleWithConstructorEmpty(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}

implementing A module app() {
	public a := ""
}
`)
}

func TestModuleWithConstructorWithArgUnused(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}

implementing A module app(str: String) {
	public a := ""
}
`)
}

func TestModuleWithConstructorWithArgUsed(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}

implementing A module app(str: String) {
	public a := str
}
`)
}

func TestModuleWithConstructorWithArgPublic(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}

implementing A module app(public a: String) {
	
}
`)
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
