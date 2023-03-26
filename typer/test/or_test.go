package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestOrVariableWithEmptyArray(t *testing.T) {
	validProgram(t, testcode.OrVariableWithEmptyArray)
}

func TestOrVariableWithTwoElementArray(t *testing.T) {
	validProgram(t, testcode.OrVariableWithTwoElementArray)
}
