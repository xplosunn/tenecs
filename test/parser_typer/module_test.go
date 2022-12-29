package parser_typer_test

import "testing"

func TestModuleEmpty(t *testing.T) {
	invalidProgram(t, `
package main

module app {
	
}
`, "module app needs to implement some interface")
}

func TestModuleWithConstructorEmpty(t *testing.T) {
	validProgram(t, `
package main

interface A {
	public a: String
}

module a(): A {
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

module a(str: String): A {
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

module a(str: String): A {
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

module a(public a: String): A {
	
}
`)
}
