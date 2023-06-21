package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestWhenExplicitExhaustive(t *testing.T) {
	validProgram(t, testcode.WhenExplicitExhaustive)
}

func TestWhenOtherSingleType(t *testing.T) {
	validProgram(t, testcode.WhenOtherSingleType)
}

func TestWhenOtherMultipleTypes(t *testing.T) {
	validProgram(t, testcode.WhenOtherMultipleTypes)
}
