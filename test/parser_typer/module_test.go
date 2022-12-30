package parser_typer_test

import "testing"

func TestModuleWithConstructorEmpty(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}

implementing A module a() {
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

implementing A module a(str: String) {
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

implementing A module a(str: String) {
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

implementing A module a(public a: String) {
	
}
`)
}
