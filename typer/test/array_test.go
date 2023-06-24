package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestArrayVariableWithEmptyArray(t *testing.T) {
	validProgram(t, testcode.ArrayVariableWithEmptyArray)
}

func TestArrayVariableWithTwoElementArray(t *testing.T) {
	validProgram(t, testcode.ArrayVariableWithTwoElementArray)
}

func TestArrayOfArray(t *testing.T) {
	validProgram(t, testcode.ArrayOfArray)
}
