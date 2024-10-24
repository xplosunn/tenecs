package parser_typer_test

import "testing"

func TestMainProgramWithEmptyListWithoutGenericAnnotated(t *testing.T) {
	invalidProgram(t, `
package main

list := []()
`, "Missing generic")
}
